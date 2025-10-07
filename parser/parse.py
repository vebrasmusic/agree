from typing import Optional, Union
import ast

import libcst as cst
import libcst.matchers as m
from libcst.display import dump
from rich import print
import time

from parser.utils import map_sqlalchemy_type


def parse_code(text: str) -> dict:
    """
    Parse Python code and extract class information.
    
    Args:
        text: Python source code as a string
        
    Returns:
        Dictionary mapping targets to classes and their fields
    """
    root = cst.parse_module(text)
    visitor = Visitor()
    root.visit(visitor)
    return visitor.index


def get_ast(text: str):
    """
    Legacy function for CLI usage with timing and printing.
    """
    start = time.perf_counter()
    index = parse_code(text)
    end = time.perf_counter()

    elapsed = (end - start) * 1000  # ms
    print(f"{elapsed:.5f} ms")
    print(index)


class Visitor(m.MatcherDecoratableVisitor):
    def __init__(self) -> None:
        super().__init__()
        self.class_call_stack: list[str] = []
        # dict [target, dict[class, some obj]]
        # pls refactor into pydantic
        self.index: dict[str, dict[str, dict[str, Union[str, int, float]]]] = {}
        self.class_dict_stack: list[dict[str, Union[str, int, float]]] = []

    def visit_ClassDef(self, node: cst.ClassDef) -> Optional[bool]:
        self.class_call_stack.append(node.name.value)
        self.class_dict_stack.append({})

    def leave_ClassDef(self, original_node: cst.ClassDef) -> None:
        target = self.class_dict_stack[-1].get("target", None)
        current_class = self.class_call_stack[-1]

        if target is None:
            return

        if self.index.get(target) is None:
            self.index[target] = {}
        self.index[target][current_class] = self.class_dict_stack[-1]
        self.class_call_stack.pop()
        self.class_dict_stack.pop()

    # onion
    #
    # per pass: create new target / class dict
    # get all class args, etc. then figure out target when we see it
    # then can push it into index as [target][class]
    #
    #

    @m.call_if_inside(m.Call(m.Name("agree")))
    @m.visit(m.Arg())
    def _get_string_args(self, node: cst.Arg) -> None:
        """
        Gets all args for the agree decorator for a class
        """

        val = None
        kw = None
        if node.keyword is not None:
            kw = node.keyword.value
        if m.matches(node.value, m.SimpleString()):
            raw = cst.ensure_type(node.value, cst.SimpleString).value
            val = ast.literal_eval(raw)  # '"event"' → 'event'
        elif m.matches(node.value, m.Integer()):
            val = cst.ensure_type(node.value, cst.Integer).value

        if val is None or kw is None:
            return

        self.class_dict_stack[-1][kw] = val

    # fire only when we're inside a ClassDef that has @agree(...)
    @m.call_if_inside(
        m.ClassDef(
            decorators=[
                m.ZeroOrMore(),
                m.Decorator(decorator=m.Call(func=m.Name("agree"))),
                m.ZeroOrMore(),
            ]
        )
    )
    def visit_Assign(self, node: cst.Assign) -> Optional[bool]:
        """
        Handle old-style SQLAlchemy Column() definitions.
        Example: id = Column(Integer, primary_key=True)
        """
        # Get the target name
        target = None
        if len(node.targets) == 1:
            assign_target = node.targets[0]
            if m.matches(assign_target.target, m.Name()):
                target = cst.ensure_type(assign_target.target, cst.Name).value

        # Skip if target is __tablename__ or similar
        if target and target.startswith("__"):
            return

        # Check if the value is a Column() call
        if not m.matches(node.value, m.Call(func=m.Name("Column"))):
            return

        call = cst.ensure_type(node.value, cst.Call)

        # Extract SQLAlchemy type and nullable from Column() arguments
        sqlalchemy_type = None
        nullable = False

        for arg in call.args:
            # First positional arg is usually the type
            if arg.keyword is None and sqlalchemy_type is None:
                if m.matches(arg.value, m.Name()):
                    sqlalchemy_type = cst.ensure_type(arg.value, cst.Name).value
            # Check for nullable keyword argument
            elif arg.keyword and m.matches(arg.keyword, m.Name(value="nullable")):
                if m.matches(arg.value, m.Name(value="True")):
                    nullable = True

        if not sqlalchemy_type or not target:
            return

        # Map SQLAlchemy type to Python type
        python_type = map_sqlalchemy_type(sqlalchemy_type)

        # Build types list
        types = [python_type]
        if nullable:
            types.append("None")

        # Store in class dict
        if "fields" not in self.class_dict_stack[-1]:
            self.class_dict_stack[-1]["fields"] = {}
        self.class_dict_stack[-1]["fields"][target] = types

    @m.call_if_inside(
        m.ClassDef(
            decorators=[
                m.ZeroOrMore(),
                m.Decorator(decorator=m.Call(func=m.Name("agree"))),
                m.ZeroOrMore(),
            ]
        )
    )
    def visit_AnnAssign(self, node: cst.AnnAssign) -> Optional[bool]:

        annotation, target = None, None

        if m.matches(node.target, m.Name()):
            target = cst.ensure_type(node.target, cst.Name).value

        # Extract all types from the annotation
        # node.annotation is an Annotation wrapper, we need the actual annotation inside
        actual_annotation = node.annotation
        if m.matches(actual_annotation, m.Annotation()):
            actual_annotation = cst.ensure_type(
                actual_annotation, cst.Annotation
            ).annotation

        annotation_types = self._extract_from_annotation(actual_annotation)

        # Store the information in the class dict
        if target and annotation_types:
            # For now, store the types list. Can be refined later based on needs
            if "fields" not in self.class_dict_stack[-1]:
                self.class_dict_stack[-1]["fields"] = {}
            self.class_dict_stack[-1]["fields"][target] = annotation_types

    def _extract_from_annotation(self, node: cst.BaseExpression) -> list[str]:
        """
        Recursively extract all type names from an annotation.
        Handles Optional[x], Union[x, y, ...], x | y, and simple types.
        Returns a normalized list of type name strings.
        """
        if m.matches(node, m.Subscript()):
            subscript = cst.ensure_type(node, cst.Subscript)
            ann_base = subscript.value

            # Get the base type name (e.g., 'Optional', 'Union', 'List')
            base_name = None
            if m.matches(ann_base, m.Name()):
                base_name = cst.ensure_type(ann_base, cst.Name).value

            # Extract types from all slice elements
            result_types = []
            for slice_element in subscript.slice:
                if m.matches(slice_element, m.SubscriptElement()):
                    element = cst.ensure_type(slice_element, cst.SubscriptElement)
                    if m.matches(element.slice, m.Index()):
                        index = cst.ensure_type(element.slice, cst.Index)
                        # Recursively extract types from the index value
                        result_types.extend(self._extract_from_annotation(index.value))

            # If it's Optional, add None to the types
            if base_name == "Optional" and "None" not in result_types:
                result_types.append("None")

            return result_types

        elif m.matches(node, m.BinaryOperation()):
            binary_op = cst.ensure_type(node, cst.BinaryOperation)
            left_base = binary_op.left
            right_base = binary_op.right
            operator = binary_op.operator

            # Check if it's the | (BitOr) operator
            if m.matches(operator, m.BitOr()):
                # Recursively extract from both sides
                left_types = self._extract_from_annotation(left_base)
                right_types = self._extract_from_annotation(right_base)
                # Combine both lists
                return left_types + right_types

            # If it's not BitOr, treat as unknown and return empty
            return []

        else:
            # Base case: simple Name node (int, str, None, etc.)
            if m.matches(node, m.Name()):
                name_node = cst.ensure_type(node, cst.Name)
                return [name_node.value]

            # For any other type we don't handle, return empty list
            return []

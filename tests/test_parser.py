"""Unit tests for parser functionality"""
import pytest
from parser.parse import parse_code


class TestPydanticSchemas:
    """Test parsing of Pydantic schema classes"""
    
    def test_simple_pydantic_schema(self):
        """Test basic Pydantic schema with simple types"""
        code = '''
from pydantic import BaseModel

@agree(target="User")
class UserSchema(BaseModel):
    id: int
    name: str
    age: float
    active: bool
'''
        result = parse_code(code)
        
        assert "User" in result
        assert "UserSchema" in result["User"]
        assert result["User"]["UserSchema"]["target"] == "User"
        assert result["User"]["UserSchema"]["fields"] == {
            "id": ["int"],
            "name": ["str"],
            "age": ["float"],
            "active": ["bool"],
        }
    
    def test_pydantic_optional_types(self):
        """Test Pydantic schema with Optional types"""
        code = '''
from typing import Optional
from pydantic import BaseModel

@agree(target="Product")
class ProductSchema(BaseModel):
    id: int
    name: str
    description: Optional[str]
    price: Optional[float]
'''
        result = parse_code(code)
        
        assert result["Product"]["ProductSchema"]["fields"] == {
            "id": ["int"],
            "name": ["str"],
            "description": ["str", "None"],
            "price": ["float", "None"],
        }
    
    def test_pydantic_union_types(self):
        """Test Pydantic schema with Union types"""
        code = '''
from typing import Union
from pydantic import BaseModel

@agree(target="Item")
class ItemSchema(BaseModel):
    id: Union[int, str]
    value: Union[int, float, str]
'''
        result = parse_code(code)
        
        assert result["Item"]["ItemSchema"]["fields"] == {
            "id": ["int", "str"],
            "value": ["int", "float", "str"],
        }
    
    def test_pydantic_pipe_unions(self):
        """Test Pydantic schema with pipe union syntax"""
        code = '''
from pydantic import BaseModel

@agree(target="Data")
class DataSchema(BaseModel):
    id: int | str
    value: int | float | None
'''
        result = parse_code(code)
        
        assert result["Data"]["DataSchema"]["fields"] == {
            "id": ["int", "str"],
            "value": ["int", "float", "None"],
        }
    
    def test_pydantic_nested_unions(self):
        """Test Pydantic schema with nested complex types"""
        code = '''
from typing import Optional, Union
from pydantic import BaseModel

@agree(target="Complex")
class ComplexSchema(BaseModel):
    nested: Optional[Union[int, str]]
    multi: Union[int, str, float, bool]
'''
        result = parse_code(code)
        
        assert result["Complex"]["ComplexSchema"]["fields"] == {
            "nested": ["int", "str", "None"],
            "multi": ["int", "str", "float", "bool"],
        }


class TestSQLAlchemyNewStyle:
    """Test parsing of new-style SQLAlchemy models with Mapped[]"""
    
    def test_simple_mapped_types(self):
        """Test basic SQLAlchemy model with Mapped types"""
        code = '''
from sqlalchemy.orm import DeclarativeBase, Mapped

class Base(DeclarativeBase):
    pass

@agree(target="User")
class UserModel(Base):
    __tablename__ = "user"
    
    id: Mapped[int]
    name: Mapped[str]
    age: Mapped[float]
'''
        result = parse_code(code)
        
        assert "User" in result
        assert result["User"]["UserModel"]["target"] == "User"
        assert result["User"]["UserModel"]["fields"] == {
            "id": ["int"],
            "name": ["str"],
            "age": ["float"],
        }
    
    def test_mapped_optional_types(self):
        """Test SQLAlchemy model with Mapped[Optional[]]"""
        code = '''
from typing import Optional
from sqlalchemy.orm import DeclarativeBase, Mapped

class Base(DeclarativeBase):
    pass

@agree(target="Product")
class ProductModel(Base):
    __tablename__ = "product"
    
    id: Mapped[int]
    name: Mapped[str]
    description: Mapped[Optional[str]]
    price: Mapped[Optional[float]]
'''
        result = parse_code(code)
        
        assert result["Product"]["ProductModel"]["fields"] == {
            "id": ["int"],
            "name": ["str"],
            "description": ["str", "None"],
            "price": ["float", "None"],
        }
    
    def test_mapped_datetime_types(self):
        """Test SQLAlchemy model with datetime types"""
        code = '''
from datetime import datetime, date
from typing import Optional
from sqlalchemy.orm import DeclarativeBase, Mapped

class Base(DeclarativeBase):
    pass

@agree(target="Event")
class EventModel(Base):
    __tablename__ = "event"
    
    id: Mapped[int]
    created_at: Mapped[datetime]
    event_date: Mapped[date]
    updated_at: Mapped[Optional[datetime]]
'''
        result = parse_code(code)
        
        assert result["Event"]["EventModel"]["fields"] == {
            "id": ["int"],
            "created_at": ["datetime"],
            "event_date": ["date"],
            "updated_at": ["datetime", "None"],
        }


class TestSQLAlchemyOldStyle:
    """Test parsing of old-style SQLAlchemy models with Column()"""
    
    def test_simple_column_types(self):
        """Test basic SQLAlchemy model with Column()"""
        code = '''
from sqlalchemy import Column, Integer, String, Float
from sqlalchemy.orm import DeclarativeBase

class Base(DeclarativeBase):
    pass

@agree(target="User")
class UserModel(Base):
    __tablename__ = "user"
    
    id = Column(Integer, primary_key=True)
    name = Column(String)
    age = Column(Float)
'''
        result = parse_code(code)
        
        assert "User" in result
        assert result["User"]["UserModel"]["target"] == "User"
        assert result["User"]["UserModel"]["fields"] == {
            "id": ["int"],
            "name": ["str"],
            "age": ["float"],
        }
    
    def test_column_nullable(self):
        """Test SQLAlchemy Column with nullable=True"""
        code = '''
from sqlalchemy import Column, Integer, String
from sqlalchemy.orm import DeclarativeBase

class Base(DeclarativeBase):
    pass

@agree(target="Product")
class ProductModel(Base):
    __tablename__ = "product"
    
    id = Column(Integer, primary_key=True)
    name = Column(String)
    description = Column(String, nullable=True)
    price = Column(Integer, nullable=True)
'''
        result = parse_code(code)
        
        assert result["Product"]["ProductModel"]["fields"] == {
            "id": ["int"],
            "name": ["str"],
            "description": ["str", "None"],
            "price": ["int", "None"],
        }
    
    def test_column_various_types(self):
        """Test SQLAlchemy Column with various type mappings"""
        code = '''
from sqlalchemy import Column, Integer, BigInteger, String, Text, Float, Boolean, DateTime, LargeBinary, JSON
from sqlalchemy.orm import DeclarativeBase

class Base(DeclarativeBase):
    pass

@agree(target="AllTypes")
class AllTypesModel(Base):
    __tablename__ = "all_types"
    
    int_field = Column(Integer)
    bigint_field = Column(BigInteger)
    string_field = Column(String)
    text_field = Column(Text)
    float_field = Column(Float)
    bool_field = Column(Boolean)
    datetime_field = Column(DateTime)
    binary_field = Column(LargeBinary)
    json_field = Column(JSON)
'''
        result = parse_code(code)
        
        assert result["AllTypes"]["AllTypesModel"]["fields"] == {
            "int_field": ["int"],
            "bigint_field": ["int"],
            "string_field": ["str"],
            "text_field": ["str"],
            "float_field": ["float"],
            "bool_field": ["bool"],
            "datetime_field": ["datetime"],
            "binary_field": ["bytes"],
            "json_field": ["dict"],
        }
    
    def test_column_skips_tablename(self):
        """Test that __tablename__ is properly skipped"""
        code = '''
from sqlalchemy import Column, Integer
from sqlalchemy.orm import DeclarativeBase

class Base(DeclarativeBase):
    pass

@agree(target="Test")
class TestModel(Base):
    __tablename__ = "test"
    
    id = Column(Integer)
'''
        result = parse_code(code)
        
        # Should only have 'id' field, not '__tablename__'
        assert "__tablename__" not in result["Test"]["TestModel"]["fields"]
        assert result["Test"]["TestModel"]["fields"] == {"id": ["int"]}


class TestMixedStyles:
    """Test that different styles can coexist and produce consistent results"""
    
    def test_same_target_multiple_classes(self):
        """Test multiple classes targeting the same entity"""
        code = '''
from typing import Optional
from sqlalchemy import Column, Integer, String
from sqlalchemy.orm import DeclarativeBase, Mapped
from pydantic import BaseModel

class Base(DeclarativeBase):
    pass

@agree(target="User")
class UserSchema(BaseModel):
    id: int
    name: str
    email: Optional[str]

@agree(target="User")
class UserModelOld(Base):
    __tablename__ = "user_old"
    id = Column(Integer)
    name = Column(String)
    email = Column(String, nullable=True)

@agree(target="User")
class UserModelNew(Base):
    __tablename__ = "user_new"
    id: Mapped[int]
    name: Mapped[str]
    email: Mapped[Optional[str]]
'''
        result = parse_code(code)
        
        # All three should be under the same target
        assert "User" in result
        assert len(result["User"]) == 3
        
        # All three should have identical field structure
        expected_fields = {
            "id": ["int"],
            "name": ["str"],
            "email": ["str", "None"],
        }
        
        assert result["User"]["UserSchema"]["fields"] == expected_fields
        assert result["User"]["UserModelOld"]["fields"] == expected_fields
        assert result["User"]["UserModelNew"]["fields"] == expected_fields
    
    def test_multiple_targets(self):
        """Test parsing multiple different targets"""
        code = '''
from pydantic import BaseModel

@agree(target="User")
class UserSchema(BaseModel):
    id: int
    name: str

@agree(target="Product")
class ProductSchema(BaseModel):
    id: int
    title: str

@agree(target="Order")
class OrderSchema(BaseModel):
    id: int
    user_id: int
'''
        result = parse_code(code)
        
        assert "User" in result
        assert "Product" in result
        assert "Order" in result
        assert len(result) == 3


class TestEdgeCases:
    """Test edge cases and error handling"""
    
    def test_class_without_agree_decorator(self):
        """Test that classes without @agree decorator are ignored"""
        code = '''
from pydantic import BaseModel

class RegularClass(BaseModel):
    id: int
    name: str

@agree(target="User")
class UserSchema(BaseModel):
    id: int
'''
        result = parse_code(code)
        
        # Should only have User, not RegularClass
        assert "User" in result
        assert len(result) == 1
    
    def test_empty_class(self):
        """Test class with no fields"""
        code = '''
from pydantic import BaseModel

@agree(target="Empty")
class EmptySchema(BaseModel):
    pass
'''
        result = parse_code(code)
        
        assert "Empty" in result
        assert result["Empty"]["EmptySchema"]["target"] == "Empty"
        # Fields key might not exist or be empty
        fields = result["Empty"]["EmptySchema"].get("fields", {})
        assert len(fields) == 0
    
    def test_large_union(self):
        """Test handling of large union types"""
        code = '''
from pydantic import BaseModel

@agree(target="BigUnion")
class BigUnionSchema(BaseModel):
    many_types: int | str | float | bool | bytes | None
'''
        result = parse_code(code)
        
        assert result["BigUnion"]["BigUnionSchema"]["fields"]["many_types"] == [
            "int", "str", "float", "bool", "bytes", "None"
        ]
    
    def test_deeply_nested_optional(self):
        """Test deeply nested Optional types"""
        code = '''
from typing import Optional
from pydantic import BaseModel

@agree(target="Nested")
class NestedSchema(BaseModel):
    triple: Optional[Optional[Optional[int]]]
'''
        result = parse_code(code)
        
        # Should normalize to just int and None
        assert result["Nested"]["NestedSchema"]["fields"]["triple"] == ["int", "None"]


class TestAgreeDecorator:
    """Test @agree decorator parameter handling"""
    
    def test_target_parameter(self):
        """Test that target parameter is correctly extracted"""
        code = '''
from pydantic import BaseModel

@agree(target="CustomTarget")
class MySchema(BaseModel):
    id: int
'''
        result = parse_code(code)
        
        assert "CustomTarget" in result
        assert result["CustomTarget"]["MySchema"]["target"] == "CustomTarget"
    
    def test_multiple_decorator_args(self):
        """Test @agree with multiple parameters"""
        code = '''
from pydantic import BaseModel

@agree(target="User", fidelity=2)
class UserSchema(BaseModel):
    id: int
'''
        result = parse_code(code)
        
        assert result["User"]["UserSchema"]["target"] == "User"
        # fidelity is stored as string representation of the integer
        assert result["User"]["UserSchema"]["fidelity"] == "2"

package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	ts "github.com/tree-sitter/go-tree-sitter"
)

// SchemaGrammar defines how to parse a specific schema format
type SchemaGrammar struct {
	Name        string            `json:"name"`
	Language    string            `json:"language"`
	Patterns    []PatternRule     `json:"patterns"`
	TypeMapping map[string]string `json:"type_mapping"`
}

// PatternRule defines a specific syntax pattern within a schema format
type PatternRule struct {
	Name       string         `json:"name"`
	Query      string         `json:"query"`
	FieldName  FieldExtractor `json:"field_name"`
	FieldType  FieldExtractor `json:"field_type"`
	Conditions []string       `json:"conditions"`
}

// FieldExtractor defines how to extract a field name or type from AST nodes
type FieldExtractor struct {
	NodeType    string `json:"node_type"`
	FieldName   string `json:"field_name"`
	ChildIndex  *int   `json:"child_index"`
	TextPattern string `json:"text_pattern"`
}

// GrammarEngine processes schema code using grammar definitions
type GrammarEngine struct {
	grammars map[string]SchemaGrammar
}

// NewGrammarEngine creates a new grammar engine
func NewGrammarEngine() *GrammarEngine {
	return &GrammarEngine{
		grammars: make(map[string]SchemaGrammar),
	}
}

// LoadGrammar loads a grammar definition from a JSON file
func (ge *GrammarEngine) LoadGrammar(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read grammar file %s: %w", filepath, err)
	}

	var grammar SchemaGrammar
	if err := json.Unmarshal(data, &grammar); err != nil {
		return fmt.Errorf("failed to parse grammar file %s: %w", filepath, err)
	}

	ge.grammars[grammar.Name] = grammar
	return nil
}

// AddGrammar adds a grammar definition directly
func (ge *GrammarEngine) AddGrammar(grammar SchemaGrammar) {
	ge.grammars[grammar.Name] = grammar
}

// ParseModel parses a model using the specified grammar
func (ge *GrammarEngine) ParseModel(src []byte, grammarName string, language *ts.Language) (Model, error) {
	grammar, exists := ge.grammars[grammarName]
	if !exists {
		return Model{}, fmt.Errorf("grammar '%s' not found", grammarName)
	}

	parser := ts.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.Parse(src, nil)
	defer tree.Close()

	root := tree.RootNode()
	return ge.extractModelFromAST(root, src, grammar)
}

// ParseTypeScriptModel parses a TypeScript model using the specified grammar
func (ge *GrammarEngine) ParseTypeScriptModel(src []byte, grammarName string, language *ts.Language) (Model, error) {
	grammar, exists := ge.grammars[grammarName]
	if !exists {
		return Model{}, fmt.Errorf("grammar '%s' not found", grammarName)
	}

	parser := ts.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.Parse(src, nil)
	defer tree.Close()

	root := tree.RootNode()
	return ge.extractTypeScriptModelFromAST(root, src, grammar)
}

// extractModelFromAST extracts a model from the AST using grammar rules
func (ge *GrammarEngine) extractModelFromAST(root *ts.Node, src []byte, grammar SchemaGrammar) (Model, error) {
	// First, find the class definition
	className := ""
	var classBody *ts.Node

	for i := uint(0); i < root.NamedChildCount(); i++ {
		n := root.NamedChild(i)
		if n.Kind() == "class_definition" {
			nameNode := n.ChildByFieldName("name")
			if nameNode != nil {
				className = nameNode.Utf8Text(src)
				classBody = n.ChildByFieldName("body")
				break
			}
		}
	}

	if classBody == nil {
		return Model{}, fmt.Errorf("no class definition found")
	}

	fields := make(map[string]Field)

	// Try each pattern in the grammar
	for _, pattern := range grammar.Patterns {
		patternFields, err := ge.extractFieldsWithPattern(classBody, src, pattern, grammar.TypeMapping)
		if err != nil {
			continue // Try next pattern
		}

		// Merge fields from this pattern
		for name, field := range patternFields {
			fields[name] = field
		}
	}

	return Model{Name: className, Fields: fields}, nil
}

// extractFieldsWithPattern extracts fields using a specific grammar pattern
func (ge *GrammarEngine) extractFieldsWithPattern(classBody *ts.Node, src []byte, pattern PatternRule, typeMapping map[string]string) (map[string]Field, error) {
	fields := make(map[string]Field)

	// Walk through all statements in the class body
	for i := uint(0); i < classBody.NamedChildCount(); i++ {
		stmt := classBody.NamedChild(i)

		// Try to match this statement with the pattern
		if field, matched := ge.matchPattern(stmt, src, pattern, typeMapping); matched {
			fields[field.Name] = field
		}
	}

	return fields, nil
}

// matchPattern tries to match a single AST node against a pattern rule
func (ge *GrammarEngine) matchPattern(node *ts.Node, src []byte, pattern PatternRule, typeMapping map[string]string) (Field, bool) {
	// For now, implement basic pattern matching for assignment statements
	// In a full implementation, this would use tree-sitter queries

	if node.Kind() != "expression_statement" {
		return Field{}, false
	}

	assign := node.NamedChild(0)
	if assign == nil || assign.Kind() != "assignment" {
		return Field{}, false
	}

	// Extract field name
	fieldName := ge.extractFieldValue(assign, src, pattern.FieldName)
	if fieldName == "" {
		return Field{}, false
	}

	// Skip dunder methods
	if strings.HasPrefix(fieldName, "__") {
		return Field{}, false
	}

	// Extract field type based on pattern
	fieldType := ge.extractFieldValue(assign, src, pattern.FieldType)

	// Apply conditions
	if !ge.checkConditions(assign, src, pattern.Conditions) {
		return Field{}, false
	}

	// Normalize type using mapping
	if mappedType, exists := typeMapping[fieldType]; exists {
		fieldType = mappedType
	}

	return Field{Name: fieldName, Type: fieldType}, true
}

// extractFieldValue extracts a value using a FieldExtractor
func (ge *GrammarEngine) extractFieldValue(node *ts.Node, src []byte, extractor FieldExtractor) string {
	var targetNode *ts.Node

	// Get the target node based on extractor config
	if extractor.FieldName != "" {
		targetNode = node.ChildByFieldName(extractor.FieldName)
	} else if extractor.ChildIndex != nil {
		if *extractor.ChildIndex < int(node.NamedChildCount()) {
			targetNode = node.NamedChild(uint(*extractor.ChildIndex))
		}
	} else {
		targetNode = node
	}

	if targetNode == nil {
		return ""
	}

	text := targetNode.Utf8Text(src)

	// Apply text pattern if specified
	if extractor.TextPattern != "" {
		re := regexp.MustCompile(extractor.TextPattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}
		return ""
	}

	return text
}

// checkConditions checks if pattern conditions are met
func (ge *GrammarEngine) checkConditions(node *ts.Node, src []byte, conditions []string) bool {
	for _, condition := range conditions {
		if !ge.evaluateCondition(node, src, condition) {
			return false
		}
	}
	return true
}

// evaluateCondition evaluates a single condition
func (ge *GrammarEngine) evaluateCondition(node *ts.Node, src []byte, condition string) bool {
	// Simple condition evaluation - can be extended
	if condition == "inside_class_body" {
		return true // We're already filtering to class body
	}

	// Handle function name conditions like "func_name == 'Column'"
	if strings.Contains(condition, "func_name ==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			expectedName := strings.Trim(strings.TrimSpace(parts[1]), "'\"")
			right := node.ChildByFieldName("right")
			if right != nil && right.Kind() == "call" {
				fn := right.ChildByFieldName("function")
				if fn != nil && fn.Utf8Text(src) == expectedName {
					return true
				}
			}
		}
	}

	return false
}

// extractTypeScriptModelFromAST extracts a TypeScript model from the AST using grammar rules
func (ge *GrammarEngine) extractTypeScriptModelFromAST(root *ts.Node, src []byte, grammar SchemaGrammar) (Model, error) {
	// For TypeScript, look for variable declarations of the form:
	// export const UserSchema = z.object({ ... })
	
	var modelName string
	var objectExpr *ts.Node

	for i := uint(0); i < root.NamedChildCount(); i++ {
		n := root.NamedChild(i)
		if n.Kind() == "export_statement" {
			// Check for variable declaration inside export
			for j := uint(0); j < n.NamedChildCount(); j++ {
				child := n.NamedChild(j)
				if child.Kind() == "lexical_declaration" {
					// Look for variable declarator
					for k := uint(0); k < child.NamedChildCount(); k++ {
						declarator := child.NamedChild(k)
						if declarator.Kind() == "variable_declarator" {
							// Get variable name
							name := declarator.ChildByFieldName("name")
							if name != nil && strings.HasSuffix(name.Utf8Text(src), "Schema") {
								modelName = strings.TrimSuffix(name.Utf8Text(src), "Schema")
								
								// Get the value (should be z.object(...))
								value := declarator.ChildByFieldName("value")
								if value != nil && value.Kind() == "call_expression" {
									// Look for the object argument
									args := value.ChildByFieldName("arguments")
									if args != nil && args.NamedChildCount() > 0 {
										firstArg := args.NamedChild(0)
										if firstArg.Kind() == "object" {
											objectExpr = firstArg
											break
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if objectExpr == nil {
		return Model{}, fmt.Errorf("no schema object found")
	}

	fields := make(map[string]Field)

	// Parse object properties
	for i := uint(0); i < objectExpr.NamedChildCount(); i++ {
		prop := objectExpr.NamedChild(i)
		if prop.Kind() == "pair" {
			// Extract field name and type
			if field, matched := ge.matchTypeScriptPattern(prop, src, grammar.Patterns[0], grammar.TypeMapping); matched {
				fields[field.Name] = field
			}
		}
	}

	return Model{Name: modelName, Fields: fields}, nil
}

// matchTypeScriptPattern matches TypeScript object property patterns
func (ge *GrammarEngine) matchTypeScriptPattern(node *ts.Node, src []byte, pattern PatternRule, typeMapping map[string]string) (Field, bool) {
	if node.Kind() != "pair" {
		return Field{}, false
	}

	// Extract field name
	key := node.ChildByFieldName("key")
	if key == nil {
		return Field{}, false
	}
	fieldName := key.Utf8Text(src)

	// Extract field type  
	value := node.ChildByFieldName("value")
	if value == nil {
		return Field{}, false
	}

	fieldType := ge.extractTypeScriptType(value, src)

	// Apply type mapping
	if mappedType, exists := typeMapping[fieldType]; exists {
		fieldType = mappedType
	}

	return Field{Name: fieldName, Type: fieldType}, true
}

// extractTypeScriptType extracts type information from TypeScript value expressions
func (ge *GrammarEngine) extractTypeScriptType(node *ts.Node, src []byte) string {
	text := node.Utf8Text(src)
	
	// Handle different Zod patterns
	if strings.Contains(text, "z.string()") {
		if strings.Contains(text, ".email()") {
			return "string().email"
		} else if strings.Contains(text, ".nullable()") {
			return "string().nullable"
		} else if strings.Contains(text, ".optional()") {
			return "string().optional"
		}
		return "string"
	} else if strings.Contains(text, "z.number()") {
		if strings.Contains(text, ".nullable()") {
			return "number().nullable"
		} else if strings.Contains(text, ".optional()") {
			return "number().optional"
		}
		return "number"
	} else if strings.Contains(text, "z.boolean()") {
		if strings.Contains(text, ".nullable()") {
			return "boolean().nullable"
		} else if strings.Contains(text, ".optional()") {
			return "boolean().optional"
		}
		return "boolean"
	} else if strings.Contains(text, "z.date()") {
		return "date"
	} else if strings.Contains(text, "z.array(") {
		if strings.Contains(text, "z.string()") {
			return "array(string())"
		} else if strings.Contains(text, "z.number()") {
			return "array(number())"
		}
		return "array"
	} else if node.Kind() == "identifier" {
		// Handle nested schema references
		return "object"
	}
	
	return "unknown"
}

// GetGrammarNames returns all loaded grammar names
func (ge *GrammarEngine) GetGrammarNames() []string {
	names := make([]string, 0, len(ge.grammars))
	for name := range ge.grammars {
		names = append(names, name)
	}
	return names
}


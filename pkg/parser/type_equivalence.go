package parser

import "strings"

// TypeEquivalenceMap defines cross-language type equivalences
type TypeEquivalenceMap struct {
	equivalences map[string][]string
}

// NewTypeEquivalenceMap creates a new type equivalence mapper
func NewTypeEquivalenceMap() *TypeEquivalenceMap {
	return &TypeEquivalenceMap{
		equivalences: map[string][]string{
			// Numeric types - TypeScript number encompasses both int and float
			"number":  {"integer", "int", "float", "number"},
			"integer": {"number", "int", "integer"},
			"int":     {"number", "integer", "int"},
			"float":   {"number", "float"},
			
			// String types
			"string": {"str", "string", "text"},
			"str":    {"string", "str", "text"},
			"text":   {"string", "str", "text"},
			
			// Boolean types
			"boolean": {"bool", "boolean"},
			"bool":    {"boolean", "bool"},
			
			// Email types (various validators)
			"email":    {"emailstr", "email", "string().email"},
			"emailstr": {"email", "emailstr", "string().email"},
			"string().email": {"email", "emailstr", "string().email"},
			
			// URL types
			"url":           {"string().url", "url"},
			"string().url":  {"url", "string().url"},
			
			// UUID types
			"uuid":          {"string().uuid", "uuid"},
			"string().uuid": {"uuid", "string().uuid"},
			
			// Date/time types
			"datetime": {"date", "datetime", "timestamp"},
			"date":     {"datetime", "date", "timestamp"},
			"timestamp": {"datetime", "date", "timestamp"},
			
			// Array types - simplified for now
			"array":     {"array", "list"},
			"list":      {"array", "list"},
			"string[]":  {"array(string())", "list[str]", "string[]"},
			"number[]":  {"array(number())", "list[int]", "list[float]", "number[]"},
			
			// Object/JSON types
			"object": {"object", "dict", "json"},
			"dict":   {"object", "dict", "json"},
			"json":   {"object", "dict", "json"},
		},
	}
}

// AreTypesEquivalent checks if two types are equivalent across languages
func (tem *TypeEquivalenceMap) AreTypesEquivalent(type1, type2 string) bool {
	// Exact match
	if type1 == type2 {
		return true
	}
	
	// Normalize types (remove whitespace, convert to lowercase)
	type1 = strings.ToLower(strings.TrimSpace(type1))
	type2 = strings.ToLower(strings.TrimSpace(type2))
	
	// Check exact match after normalization
	if type1 == type2 {
		return true
	}
	
	// Handle nullable types (optional in some languages)
	type1Clean, type1Nullable := extractNullableType(type1)
	type2Clean, type2Nullable := extractNullableType(type2)
	
	// If one is nullable and the other isn't, they're not equivalent
	if type1Nullable != type2Nullable {
		return false
	}
	
	// Check equivalence of the base types
	return tem.areBaseTypesEquivalent(type1Clean, type2Clean)
}

// areBaseTypesEquivalent checks if base types (without nullable modifiers) are equivalent
func (tem *TypeEquivalenceMap) areBaseTypesEquivalent(type1, type2 string) bool {
	// Check if type1 has equivalents that include type2
	if equivalents, exists := tem.equivalences[type1]; exists {
		for _, equiv := range equivalents {
			if equiv == type2 {
				return true
			}
		}
	}
	
	// Check if type2 has equivalents that include type1
	if equivalents, exists := tem.equivalences[type2]; exists {
		for _, equiv := range equivalents {
			if equiv == type1 {
				return true
			}
		}
	}
	
	return false
}

// extractNullableType extracts the base type and nullable status
func extractNullableType(typeStr string) (baseType string, isNullable bool) {
	typeStr = strings.TrimSpace(typeStr)
	
	// Handle different nullable patterns
	if strings.HasSuffix(typeStr, "?") {
		return strings.TrimSuffix(typeStr, "?"), true
	}
	
	if strings.HasSuffix(typeStr, "| none") {
		return strings.TrimSpace(strings.TrimSuffix(typeStr, "| none")), true
	}
	
	if strings.HasSuffix(typeStr, "| null") {
		return strings.TrimSpace(strings.TrimSuffix(typeStr, "| null")), true
	}
	
	if strings.Contains(typeStr, "optional[") {
		// Extract Optional[Type] -> Type
		start := strings.Index(typeStr, "[")
		end := strings.LastIndex(typeStr, "]")
		if start != -1 && end != -1 && end > start {
			return strings.TrimSpace(typeStr[start+1:end]), true
		}
	}
	
	if strings.Contains(typeStr, ".nullable") {
		// Handle z.string().nullable() -> string
		if strings.Contains(typeStr, "string") {
			return "string", true
		}
		if strings.Contains(typeStr, "number") {
			return "number", true
		}
		if strings.Contains(typeStr, "boolean") {
			return "boolean", true
		}
	}
	
	if strings.Contains(typeStr, ".optional") {
		// Handle z.string().optional() -> string
		if strings.Contains(typeStr, "string") {
			return "string", true
		}
		if strings.Contains(typeStr, "number") {
			return "number", true
		}
		if strings.Contains(typeStr, "boolean") {
			return "boolean", true
		}
	}
	
	return typeStr, false
}

// GetCanonicalType returns a canonical type representation for comparison
func (tem *TypeEquivalenceMap) GetCanonicalType(typeStr string) string {
	baseType, isNullable := extractNullableType(strings.ToLower(strings.TrimSpace(typeStr)))
	
	// Map to canonical types
	canonical := baseType
	switch baseType {
	case "int", "integer":
		canonical = "integer"
	case "str":
		canonical = "string"
	case "bool":
		canonical = "boolean"
	case "emailstr", "string().email":
		canonical = "email"
	case "datetime", "timestamp":
		canonical = "date"
	case "dict":
		canonical = "object"
	case "list":
		canonical = "array"
	}
	
	if isNullable {
		canonical += "?"
	}
	
	return canonical
}

// AddEquivalence adds a new type equivalence
func (tem *TypeEquivalenceMap) AddEquivalence(primaryType string, equivalentTypes ...string) {
	primaryType = strings.ToLower(strings.TrimSpace(primaryType))
	
	// Add primary type to its own equivalents
	if _, exists := tem.equivalences[primaryType]; !exists {
		tem.equivalences[primaryType] = []string{primaryType}
	}
	
	// Add all equivalent types
	for _, equivType := range equivalentTypes {
		equivType = strings.ToLower(strings.TrimSpace(equivType))
		tem.equivalences[primaryType] = append(tem.equivalences[primaryType], equivType)
		
		// Also add reverse mapping
		if _, exists := tem.equivalences[equivType]; !exists {
			tem.equivalences[equivType] = []string{equivType}
		}
		tem.equivalences[equivType] = append(tem.equivalences[equivType], primaryType)
	}
}
package parser

import (
	"testing"
)

func TestTypeEquivalence_Basic(t *testing.T) {
	tem := NewTypeEquivalenceMap()
	
	tests := []struct {
		name     string
		type1    string
		type2    string
		expected bool
	}{
		// Exact matches
		{"exact string match", "string", "string", true},
		{"exact number match", "number", "number", true},
		
		// Cross-language numeric equivalence
		{"number equals integer", "number", "integer", true},
		{"integer equals number", "integer", "number", true},
		{"number equals int", "number", "int", true},
		{"int equals number", "int", "number", true},
		{"integer equals int", "integer", "int", true},
		{"float equals number", "float", "number", true},
		
		// String equivalence
		{"string equals str", "string", "str", true},
		{"str equals string", "str", "string", true},
		{"text equals string", "text", "string", true},
		
		// Boolean equivalence
		{"boolean equals bool", "boolean", "bool", true},
		{"bool equals boolean", "bool", "boolean", true},
		
		// Email equivalence
		{"email equals emailstr", "email", "emailstr", true},
		{"emailstr equals email", "emailstr", "email", true},
		{"email equals string().email", "email", "string().email", true},
		
		// Non-equivalent types
		{"string not equal to number", "string", "number", false},
		{"boolean not equal to string", "boolean", "string", false},
		{"integer not equal to string", "integer", "string", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tem.AreTypesEquivalent(tt.type1, tt.type2)
			if result != tt.expected {
				t.Errorf("AreTypesEquivalent(%q, %q) = %v, expected %v", 
					tt.type1, tt.type2, result, tt.expected)
			}
		})
	}
}

func TestTypeEquivalence_Nullable(t *testing.T) {
	tem := NewTypeEquivalenceMap()
	
	tests := []struct {
		name     string
		type1    string
		type2    string
		expected bool
	}{
		// Nullable equivalence
		{"nullable string match", "string?", "string?", true},
		{"nullable number equals nullable integer", "number?", "integer?", true},
		
		// Nullable vs non-nullable (should not match)
		{"nullable vs non-nullable string", "string?", "string", false},
		{"nullable vs non-nullable number", "number?", "number", false},
		
		// Different nullable syntax
		{"optional syntax match", "string?", "string | none", true},
		{"nullable zod syntax", "string().nullable", "string?", true},
		{"optional zod syntax", "string().optional", "string?", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tem.AreTypesEquivalent(tt.type1, tt.type2)
			if result != tt.expected {
				t.Errorf("AreTypesEquivalent(%q, %q) = %v, expected %v", 
					tt.type1, tt.type2, result, tt.expected)
			}
		})
	}
}

func TestTypeEquivalence_CaseInsensitive(t *testing.T) {
	tem := NewTypeEquivalenceMap()
	
	tests := []struct {
		name     string
		type1    string
		type2    string
		expected bool
	}{
		{"mixed case string", "String", "string", true},
		{"mixed case number", "Number", "INTEGER", true},
		{"mixed case boolean", "Boolean", "BOOL", true},
		{"whitespace handling", " string ", "string", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tem.AreTypesEquivalent(tt.type1, tt.type2)
			if result != tt.expected {
				t.Errorf("AreTypesEquivalent(%q, %q) = %v, expected %v", 
					tt.type1, tt.type2, result, tt.expected)
			}
		})
	}
}

func TestTypeEquivalence_CrossLanguageRealistic(t *testing.T) {
	tem := NewTypeEquivalenceMap()
	
	// Test realistic cross-language scenarios
	tests := []struct {
		name        string
		pydanticType string
		zodType      string
		expected     bool
	}{
		{"id field: int vs number", "int", "number", true},
		{"price field: float vs number", "float", "number", true},
		{"name field: str vs string", "str", "string", true},
		{"active field: bool vs boolean", "bool", "boolean", true},
		{"email field: EmailStr vs email", "EmailStr", "email", true},
		{"optional name: str? vs string?", "str | None", "string?", true},
		{"zod nullable: str vs string().nullable", "str", "string().nullable", false}, // nullable mismatch
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tem.AreTypesEquivalent(tt.pydanticType, tt.zodType)
			if result != tt.expected {
				t.Errorf("Cross-language types %q and %q: got %v, expected %v", 
					tt.pydanticType, tt.zodType, result, tt.expected)
			}
		})
	}
}

func TestExtractNullableType(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedBase   string
		expectedNull   bool
	}{
		{"simple string", "string", "string", false},
		{"nullable with ?", "string?", "string", true},
		{"nullable with | none", "string | none", "string", true},
		{"nullable with | null", "string | null", "string", true},
		{"zod nullable", "string().nullable", "string", true},
		{"zod optional", "number().optional", "number", true},
		{"optional type", "optional[str]", "str", true}, // lowercase for regex match
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, nullable := extractNullableType(tt.input)
			if base != tt.expectedBase || nullable != tt.expectedNull {
				t.Errorf("extractNullableType(%q) = (%q, %v), expected (%q, %v)", 
					tt.input, base, nullable, tt.expectedBase, tt.expectedNull)
			}
		})
	}
}

func TestGetCanonicalType(t *testing.T) {
	tem := NewTypeEquivalenceMap()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normalize int", "int", "integer"},
		{"normalize str", "str", "string"},
		{"normalize bool", "bool", "boolean"},
		{"normalize nullable", "str?", "string?"},
		{"normalize EmailStr", "EmailStr", "email"},
		{"keep number", "number", "number"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tem.GetCanonicalType(tt.input)
			if result != tt.expected {
				t.Errorf("GetCanonicalType(%q) = %q, expected %q", 
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestCompareModelsWithEquivalence(t *testing.T) {
	// Test that the new comparison function works correctly
	pydanticModel := Model{
		Name: "User",
		Fields: map[string]Field{
			"id":       {Name: "id", Type: "int"},           // Should match 'number'
			"name":     {Name: "name", Type: "str"},         // Should match 'string'
			"active":   {Name: "active", Type: "bool"},      // Should match 'boolean'
			"email":    {Name: "email", Type: "EmailStr"},   // Should match 'email'
		},
	}
	
	zodModel := Model{
		Name: "User",
		Fields: map[string]Field{
			"id":     {Name: "id", Type: "number"},     // Should match 'int'
			"name":   {Name: "name", Type: "string"},   // Should match 'str'
			"active": {Name: "active", Type: "boolean"}, // Should match 'bool'
			"email":  {Name: "email", Type: "email"},    // Should match 'EmailStr'
		},
	}
	
	models1 := map[string]Model{"user": pydanticModel}
	models2 := map[string]Model{"user": zodModel}
	
	result := CompareModelsWithEquivalence(models1, models2)
	
	// Should find no mismatches due to type equivalence
	if result != "No mismatches found" {
		t.Errorf("Expected no mismatches with type equivalence, got: %s", result)
	}
}
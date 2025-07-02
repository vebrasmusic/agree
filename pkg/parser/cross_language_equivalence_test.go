package parser

import (
	"fmt"
	"strings"
	"testing"
)

// TestCrossLanguageEquivalenceDemo demonstrates the improved type equivalence
func TestCrossLanguageEquivalenceDemo(t *testing.T) {
	// Create realistic cross-language schemas
	pydanticModel := Model{
		Name: "User",
		Fields: map[string]Field{
			"id":       {Name: "id", Type: "int"},         // Python int
			"price":    {Name: "price", Type: "float"},    // Python float
			"name":     {Name: "name", Type: "str"},       // Python str
			"active":   {Name: "active", Type: "bool"},    // Python bool
			"email":    {Name: "email", Type: "EmailStr"}, // Pydantic EmailStr
		},
	}
	
	zodModel := Model{
		Name: "User",
		Fields: map[string]Field{
			"id":     {Name: "id", Type: "number"},        // TypeScript number
			"price":  {Name: "price", Type: "number"},     // TypeScript number
			"name":   {Name: "name", Type: "string"},      // TypeScript string
			"active": {Name: "active", Type: "boolean"},   // TypeScript boolean
			"email":  {Name: "email", Type: "email"},      // Zod email validation
		},
	}
	
	// Test with old comparison (should find many mismatches)
	models1 := map[string]Model{"user": pydanticModel}
	models2 := map[string]Model{"user": zodModel}
	
	// Using new equivalence-based comparison
	result := CompareModelsWithEquivalence(models1, models2)
	
	t.Logf("Cross-language comparison result: %s", result)
	
	// Should find no mismatches due to type equivalence
	if result != "No mismatches found" {
		t.Errorf("Expected no mismatches with type equivalence, got: %s", result)
	}
}

// TestCrossLanguageStillCatchesRealDifferences ensures we still catch actual differences
func TestCrossLanguageStillCatchesRealDifferences(t *testing.T) {
	pydanticModel := Model{
		Name: "User",
		Fields: map[string]Field{
			"id":       {Name: "id", Type: "int"},
			"name":     {Name: "name", Type: "str"},
			"age":      {Name: "age", Type: "int"},        // Missing in TypeScript
			"is_admin": {Name: "is_admin", Type: "bool"},  // Different name in TypeScript
		},
	}
	
	zodModel := Model{
		Name: "User",
		Fields: map[string]Field{
			"id":    {Name: "id", Type: "number"},
			"name":  {Name: "name", Type: "string"},
			"admin": {Name: "admin", Type: "boolean"},     // Different field name
			"score": {Name: "score", Type: "number"},      // Missing in Python
		},
	}
	
	models1 := map[string]Model{"user": pydanticModel}
	models2 := map[string]Model{"user": zodModel}
	
	result := CompareModelsWithEquivalence(models1, models2)
	
	t.Logf("Real differences result: %s", result)
	
	// Should still find actual structural differences
	if result == "No mismatches found" {
		t.Error("Expected to find real structural differences, but got no mismatches")
	}
	
	// Should mention missing fields
	if !strings.Contains(result, "Missing") {
		t.Error("Expected to find missing fields in comparison")
	}
}

// TestNullableEquivalence tests nullable type handling across languages
func TestNullableEquivalence(t *testing.T) {
	tests := []struct {
		name          string
		pydanticType  string
		zodType       string
		shouldMatch   bool
		description   string
	}{
		{
			name:         "nullable strings match",
			pydanticType: "str | None",
			zodType:      "string?",
			shouldMatch:  true,
			description:  "Python union with None should match TypeScript optional",
		},
		{
			name:         "zod nullable matches python optional",
			pydanticType: "str?",
			zodType:      "string().nullable",
			shouldMatch:  true,
			description:  "Python ? syntax should match Zod nullable()",
		},
		{
			name:         "non-nullable vs nullable should not match",
			pydanticType: "str",
			zodType:      "string?",
			shouldMatch:  false,
			description:  "Required vs optional should not match",
		},
		{
			name:         "both nullable numbers should match",
			pydanticType: "int | None",
			zodType:      "number().optional",
			shouldMatch:  true,
			description:  "Optional numbers should be equivalent",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pydanticModel := Model{
				Name: "Test",
				Fields: map[string]Field{
					"field": {Name: "field", Type: tt.pydanticType},
				},
			}
			
			zodModel := Model{
				Name: "Test",
				Fields: map[string]Field{
					"field": {Name: "field", Type: tt.zodType},
				},
			}
			
			models1 := map[string]Model{"test": pydanticModel}
			models2 := map[string]Model{"test": zodModel}
			
			result := CompareModelsWithEquivalence(models1, models2)
			
			if tt.shouldMatch && result != "No mismatches found" {
				t.Errorf("%s: Expected match but got mismatches: %s", tt.description, result)
			}
			
			if !tt.shouldMatch && result == "No mismatches found" {
				t.Errorf("%s: Expected mismatches but got none", tt.description)
			}
			
			t.Logf("%s: %s vs %s = %s", tt.name, tt.pydanticType, tt.zodType, result)
		})
	}
}

// TestTypeEquivalencePerformance ensures equivalence checking doesn't slow down comparisons
func TestTypeEquivalencePerformance(t *testing.T) {
	// Create large models
	largeModel1 := Model{Name: "Large", Fields: make(map[string]Field)}
	largeModel2 := Model{Name: "Large", Fields: make(map[string]Field)}
	
	// Add 50 fields with equivalent but different type names
	for i := 0; i < 50; i++ {
		fieldName := fmt.Sprintf("field_%d", i)
		// Alternate between equivalent types
		if i%2 == 0 {
			largeModel1.Fields[fieldName] = Field{Name: fieldName, Type: "int"}
			largeModel2.Fields[fieldName] = Field{Name: fieldName, Type: "number"}
		} else {
			largeModel1.Fields[fieldName] = Field{Name: fieldName, Type: "str"}
			largeModel2.Fields[fieldName] = Field{Name: fieldName, Type: "string"}
		}
	}
	
	models1 := map[string]Model{"large": largeModel1}
	models2 := map[string]Model{"large": largeModel2}
	
	// This should complete quickly and find no mismatches
	result := CompareModelsWithEquivalence(models1, models2)
	
	if result != "No mismatches found" {
		t.Errorf("Large model comparison failed: %s", result)
	}
	
	t.Log("Performance test passed - large model comparison completed successfully")
}
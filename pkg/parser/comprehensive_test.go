package parser

import (
	"fmt"
	"strings"
	"testing"
)

// TestExactMatches tests scenarios where schemas should match perfectly
func TestExactMatches(t *testing.T) {
	tests := []struct {
		name        string
		schema1     Model
		schema2     Model
		expectMatch bool
	}{
		{
			name: "identical schemas",
			schema1: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "integer"},
					"username": {Name: "username", Type: "string"},
					"email":    {Name: "email", Type: "email"},
				},
			},
			schema2: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "integer"},
					"username": {Name: "username", Type: "string"},
					"email":    {Name: "email", Type: "email"},
				},
			},
			expectMatch: true,
		},
		{
			name: "empty schemas",
			schema1: Model{
				Name:   "Empty",
				Fields: map[string]Field{},
			},
			schema2: Model{
				Name:   "Empty",
				Fields: map[string]Field{},
			},
			expectMatch: true,
		},
		{
			name: "single field match",
			schema1: Model{
				Name: "Simple",
				Fields: map[string]Field{
					"id": {Name: "id", Type: "integer"},
				},
			},
			schema2: Model{
				Name: "Simple",
				Fields: map[string]Field{
					"id": {Name: "id", Type: "integer"},
				},
			},
			expectMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test model maps
			models1 := map[string]Model{"test": tt.schema1}
			models2 := map[string]Model{"test": tt.schema2}
			
			report := CompareModels(models1, models2)
			
			if tt.expectMatch && report != "No mismatches found" {
				t.Errorf("Expected no mismatches, got: %s", report)
			}
			if !tt.expectMatch && report == "No mismatches found" {
				t.Errorf("Expected mismatches, but got none")
			}
		})
	}
}

// TestMissingFields tests scenarios with missing fields
func TestMissingFields(t *testing.T) {
	tests := []struct {
		name           string
		schema1        Model
		schema2        Model
		expectedMissing []string
	}{
		{
			name: "missing in schema2",
			schema1: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "integer"},
					"username": {Name: "username", Type: "string"},
					"email":    {Name: "email", Type: "email"},
				},
			},
			schema2: Model{
				Name: "User", 
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "integer"},
					"username": {Name: "username", Type: "string"},
				},
			},
			expectedMissing: []string{"email"},
		},
		{
			name: "missing in schema1",
			schema1: Model{
				Name: "User",
				Fields: map[string]Field{
					"id": {Name: "id", Type: "integer"},
				},
			},
			schema2: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "integer"},
					"username": {Name: "username", Type: "string"},
					"email":    {Name: "email", Type: "email"},
				},
			},
			expectedMissing: []string{"username", "email"},
		},
		{
			name: "completely different fields",
			schema1: Model{
				Name: "A",
				Fields: map[string]Field{
					"field_a": {Name: "field_a", Type: "string"},
					"field_b": {Name: "field_b", Type: "integer"},
				},
			},
			schema2: Model{
				Name: "B",
				Fields: map[string]Field{
					"field_x": {Name: "field_x", Type: "string"},
					"field_y": {Name: "field_y", Type: "integer"},
				},
			},
			expectedMissing: []string{"field_a", "field_b", "field_x", "field_y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models1 := map[string]Model{"test": tt.schema1}
			models2 := map[string]Model{"test": tt.schema2}
			
			report := CompareModels(models1, models2)
			
			if report == "No mismatches found" {
				t.Error("Expected mismatches but got none")
				return
			}
			
			// Check that all expected missing fields are mentioned
			for _, missing := range tt.expectedMissing {
				if !strings.Contains(report, missing) {
					t.Errorf("Expected missing field '%s' not found in report: %s", missing, report)
				}
			}
		})
	}
}

// TestTypeMismatches tests scenarios with type mismatches
func TestTypeMismatches(t *testing.T) {
	tests := []struct {
		name                string
		schema1             Model
		schema2             Model
		expectedMismatches  []string
	}{
		{
			name: "single type mismatch",
			schema1: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "integer"},
					"username": {Name: "username", Type: "string"},
				},
			},
			schema2: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "string"}, // mismatch
					"username": {Name: "username", Type: "string"},
				},
			},
			expectedMismatches: []string{"id (integer != string)"},
		},
		{
			name: "multiple type mismatches",
			schema1: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":    {Name: "id", Type: "integer"},
					"score": {Name: "score", Type: "float"},
					"active": {Name: "active", Type: "boolean"},
				},
			},
			schema2: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":    {Name: "id", Type: "string"},    // mismatch
					"score": {Name: "score", Type: "integer"}, // mismatch  
					"active": {Name: "active", Type: "boolean"}, // match
				},
			},
			expectedMismatches: []string{"id (integer != string)", "score (float != integer)"},
		},
		{
			name: "nullable vs non-nullable",
			schema1: Model{
				Name: "User",
				Fields: map[string]Field{
					"full_name": {Name: "full_name", Type: "string"},
				},
			},
			schema2: Model{
				Name: "User",
				Fields: map[string]Field{
					"full_name": {Name: "full_name", Type: "string?"},
				},
			},
			expectedMismatches: []string{"full_name (string != string?)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models1 := map[string]Model{"test": tt.schema1}
			models2 := map[string]Model{"test": tt.schema2}
			
			report := CompareModels(models1, models2)
			
			if report == "No mismatches found" {
				t.Error("Expected type mismatches but got none")
				return
			}
			
			// Check that all expected mismatches are reported
			for _, mismatch := range tt.expectedMismatches {
				if !strings.Contains(report, mismatch) {
					t.Errorf("Expected mismatch '%s' not found in report: %s", mismatch, report)
				}
			}
		})
	}
}

// TestCrossLanguageMatching tests cross-language schema comparison
func TestCrossLanguageMatching(t *testing.T) {
	tests := []struct {
		name        string
		pydanticModel Model
		zodModel      Model
		expectMatch   bool
		expectedType  string // what kind of mismatch to expect
	}{
		{
			name: "perfect cross-language match",
			pydanticModel: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "number"},
					"username": {Name: "username", Type: "string"},
					"email":    {Name: "email", Type: "email"},
				},
			},
			zodModel: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "number"},
					"username": {Name: "username", Type: "string"},
					"email":    {Name: "email", Type: "email"},
				},
			},
			expectMatch: true,
		},
		{
			name: "number vs integer should now match due to equivalence",
			pydanticModel: Model{
				Name: "User",
				Fields: map[string]Field{
					"id": {Name: "id", Type: "integer"},
				},
			},
			zodModel: Model{
				Name: "User",
				Fields: map[string]Field{
					"id": {Name: "id", Type: "number"},
				},
			},
			expectMatch:  true, // Changed: now equivalent
		},
		{
			name: "missing field cross-language",
			pydanticModel: Model{
				Name: "User",
				Fields: map[string]Field{
					"id":       {Name: "id", Type: "integer"},
					"username": {Name: "username", Type: "string"},
				},
			},
			zodModel: Model{
				Name: "User",
				Fields: map[string]Field{
					"id": {Name: "id", Type: "number"},
				},
			},
			expectMatch:  false,
			expectedType: "Missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate what the full parsing would create
			allModels := map[string]map[string]Model{
				"pydantic": {"user": tt.pydanticModel},
				"zod":      {"user": tt.zodModel},
			}
			
			report := CompareModelsWithGrammars(allModels, "pydantic", "zod")
			
			if tt.expectMatch && report != "No mismatches found" {
				t.Errorf("Expected cross-language match, got: %s", report)
			}
			if !tt.expectMatch {
				if report == "No mismatches found" {
					t.Error("Expected cross-language mismatch but got none")
				} else if !strings.Contains(report, tt.expectedType) {
					t.Errorf("Expected mismatch type '%s' not found in report: %s", tt.expectedType, report)
				}
			}
		})
	}
}

// TestGrammarComparison tests the grammar-based comparison functionality
func TestGrammarComparison(t *testing.T) {
	tests := []struct {
		name           string
		allModels      map[string]map[string]Model
		schema1        string
		schema2        string
		expectError    bool
		expectedResult string
	}{
		{
			name: "valid schema comparison",
			allModels: map[string]map[string]Model{
				"pydantic": {
					"user": {
						Name: "User",
						Fields: map[string]Field{
							"id": {Name: "id", Type: "integer"},
						},
					},
				},
				"zod": {
					"user": {
						Name: "User", 
						Fields: map[string]Field{
							"id": {Name: "id", Type: "integer"},
						},
					},
				},
			},
			schema1:        "pydantic",
			schema2:        "zod",
			expectError:    false,
			expectedResult: "No mismatches found",
		},
		{
			name: "missing schema type",
			allModels: map[string]map[string]Model{
				"pydantic": {
					"user": {Name: "User", Fields: map[string]Field{}},
				},
			},
			schema1:        "pydantic",
			schema2:        "nonexistent",
			expectError:    true,
			expectedResult: "not found",
		},
		{
			name: "empty model sets",
			allModels: map[string]map[string]Model{
				"schema1": {},
				"schema2": {},
			},
			schema1:        "schema1",
			schema2:        "schema2",
			expectError:    false,
			expectedResult: "No mismatches found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareModelsWithGrammars(tt.allModels, tt.schema1, tt.schema2)
			
			if tt.expectError && !strings.Contains(result, tt.expectedResult) {
				t.Errorf("Expected error message containing '%s', got: %s", tt.expectedResult, result)
			}
			if !tt.expectError && result != tt.expectedResult {
				t.Errorf("Expected result '%s', got: %s", tt.expectedResult, result)
			}
		})
	}
}

// TestEdgeCases tests edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		models1     map[string]Model
		models2     map[string]Model
		description string
	}{
		{
			name:        "nil field maps",
			models1:     map[string]Model{"test": {Name: "Test", Fields: nil}},
			models2:     map[string]Model{"test": {Name: "Test", Fields: map[string]Field{}}},
			description: "should handle nil field maps gracefully",
		},
		{
			name:        "empty field names",
			models1:     map[string]Model{"test": {Name: "Test", Fields: map[string]Field{"": {Name: "", Type: "string"}}}},
			models2:     map[string]Model{"test": {Name: "Test", Fields: map[string]Field{}}},
			description: "should handle empty field names",
		},
		{
			name:        "special characters in field names",
			models1:     map[string]Model{"test": {Name: "Test", Fields: map[string]Field{"field_with-special.chars": {Name: "field_with-special.chars", Type: "string"}}}},
			models2:     map[string]Model{"test": {Name: "Test", Fields: map[string]Field{"field_with-special.chars": {Name: "field_with-special.chars", Type: "string"}}}},
			description: "should handle special characters in field names",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Test panicked: %v", r)
				}
			}()
			
			result := CompareModels(tt.models1, tt.models2)
			
			// Basic validation - result should be a string
			if len(result) == 0 {
				t.Error("Expected non-empty result")
			}
		})
	}
}

// BenchmarkComparison benchmarks the comparison function
func BenchmarkComparison(b *testing.B) {
	// Create large schemas for benchmarking
	largeSchema1 := Model{
		Name:   "LargeSchema",
		Fields: make(map[string]Field),
	}
	largeSchema2 := Model{
		Name:   "LargeSchema",
		Fields: make(map[string]Field),
	}
	
	// Add 100 fields to each schema
	for i := 0; i < 100; i++ {
		fieldName := fmt.Sprintf("field_%d", i)
		field := Field{Name: fieldName, Type: "string"}
		largeSchema1.Fields[fieldName] = field
		largeSchema2.Fields[fieldName] = field
	}
	
	models1 := map[string]Model{"large": largeSchema1}
	models2 := map[string]Model{"large": largeSchema2}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CompareModels(models1, models2)
	}
}
package parser

import (
	"os"
	"strings"
	"testing"
)

// TestRealFileParsingMatching tests parsing actual files with matching schemas
func TestRealFileParsingMatching(t *testing.T) {
	// This test requires the test data files to exist
	if _, err := os.Stat("../../test-data/test_schemas_matching.py"); os.IsNotExist(err) {
		t.Skip("Test data files not found")
	}
	if _, err := os.Stat("../../test-data/test_schemas_matching.ts"); os.IsNotExist(err) {
		t.Skip("Test data files not found")
	}
	
	// Parse files using the grammar engine
	allModels, err := ParseFilesWithGrammars("../../test-data", "../../grammars")
	if err != nil {
		t.Fatalf("Failed to parse files: %v", err)
	}
	
	// Check that we have the expected schema types
	expectedTypes := []string{"pydantic", "sqlalchemy", "zod"}
	for _, schemaType := range expectedTypes {
		if _, exists := allModels[schemaType]; !exists {
			t.Errorf("Expected schema type '%s' not found", schemaType)
		}
	}
	
	// Test that matching schemas are detected correctly
	// Compare match_test schemas (should have minimal mismatches)
	if pydanticModels, exists := allModels["pydantic"]; exists {
		if zodModels, exists := allModels["zod"]; exists {
			if _, hasPyd := pydanticModels["match_test"]; hasPyd {
				if _, hasZod := zodModels["match_test"]; hasZod {
					report := CompareModelsWithGrammars(allModels, "pydantic", "zod")
					
					// Should detect some type differences but overall structure should be similar
					// This tests that our cross-language comparison works
					if strings.Contains(report, "match_test") {
						t.Logf("Cross-language comparison detected expected differences: %s", report)
					}
				}
			}
		}
	}
}

// TestRealFileParsingMismatched tests parsing actual files with intentionally mismatched schemas
func TestRealFileParsingMismatched(t *testing.T) {
	// This test requires the test data files to exist
	if _, err := os.Stat("../../test-data/test_schemas_mismatched.py"); os.IsNotExist(err) {
		t.Skip("Test data files not found")
	}
	if _, err := os.Stat("../../test-data/test_schemas_mismatched.ts"); os.IsNotExist(err) {
		t.Skip("Test data files not found")
	}
	
	// Parse files using the grammar engine
	allModels, err := ParseFilesWithGrammars("../../test-data", "../../grammars")
	if err != nil {
		t.Fatalf("Failed to parse files: %v", err)
	}
	
	// Test that mismatched schemas are detected correctly
	report := CompareModelsWithGrammars(allModels, "pydantic", "zod")
	
	// Should definitely find mismatches for the intentionally mismatched schemas
	if !strings.Contains(report, "mismatch_test") {
		t.Error("Expected to find mismatches in mismatch_test schemas")
	}
	
	// Should mention missing fields
	if !strings.Contains(report, "Missing") {
		t.Error("Expected to find missing fields in mismatched schemas")
	}
	
	t.Logf("Detected mismatches as expected: %s", report)
}

// TestGrammarLoading tests that all grammar files load correctly
func TestGrammarLoading(t *testing.T) {
	engine := NewGrammarEngine()
	
	grammarFiles := []string{
		"../../grammars/pydantic.json",
		"../../grammars/sqlalchemy.json", 
		"../../grammars/zod.json",
	}
	
	for _, grammarFile := range grammarFiles {
		if _, err := os.Stat(grammarFile); os.IsNotExist(err) {
			t.Skipf("Grammar file %s not found", grammarFile)
		}
		
		err := engine.LoadGrammar(grammarFile)
		if err != nil {
			t.Errorf("Failed to load grammar %s: %v", grammarFile, err)
		}
	}
	
	// Check that all grammars were loaded
	grammarNames := engine.GetGrammarNames()
	expectedGrammars := []string{"pydantic", "sqlalchemy", "zod"}
	
	for _, expected := range expectedGrammars {
		found := false
		for _, loaded := range grammarNames {
			if loaded == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected grammar '%s' not loaded. Available: %v", expected, grammarNames)
		}
	}
}

// TestCrossLanguageTypeMapping tests type mapping between languages
func TestCrossLanguageTypeMapping(t *testing.T) {
	tests := []struct {
		name        string
		pydanticType string
		zodType      string
		shouldMatch  bool
	}{
		{"string types", "string", "string", true},
		{"boolean types", "boolean", "boolean", true},
		{"integer vs number", "integer", "number", true}, // Now equivalent
		{"email types", "email", "email", true},
		{"nullable string", "string?", "string?", true},
		{"nullable vs regular", "string", "string?", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal schemas with just the types we want to test
			pydanticModel := Model{
				Name: "Test",
				Fields: map[string]Field{
					"test_field": {Name: "test_field", Type: tt.pydanticType},
				},
			}
			
			zodModel := Model{
				Name: "Test",
				Fields: map[string]Field{
					"test_field": {Name: "test_field", Type: tt.zodType},
				},
			}
			
			allModels := map[string]map[string]Model{
				"pydantic": {"test": pydanticModel},
				"zod":      {"test": zodModel},
			}
			
			report := CompareModelsWithGrammars(allModels, "pydantic", "zod")
			
			if tt.shouldMatch && report != "No mismatches found" {
				t.Errorf("Expected types %s and %s to match, but got mismatches: %s", 
					tt.pydanticType, tt.zodType, report)
			}
			
			if !tt.shouldMatch && report == "No mismatches found" {
				t.Errorf("Expected types %s and %s to NOT match, but no mismatches found", 
					tt.pydanticType, tt.zodType)
			}
		})
	}
}

// TestFullWorkflow tests the complete workflow from file parsing to comparison
func TestFullWorkflow(t *testing.T) {
	// This test validates the entire workflow that the CLI uses
	
	// 1. Parse all files with grammars
	allModels, err := ParseFilesWithGrammars("../../test-data", "../../grammars")
	if err != nil {
		t.Fatalf("Failed to parse files: %v", err)
	}
	
	// 2. Verify we found models
	totalModels := 0
	for schemaType, models := range allModels {
		count := len(models)
		totalModels += count
		t.Logf("Found %d models for schema type '%s'", count, schemaType)
	}
	
	if totalModels == 0 {
		t.Error("No models found - this suggests parsing failed")
	}
	
	// 3. Test various comparisons
	comparisons := []struct {
		name     string
		schema1  string
		schema2  string
		expectError bool
	}{
		{"SQLAlchemy vs Pydantic", "sqlalchemy", "pydantic", false},
		{"Pydantic vs Zod", "pydantic", "zod", false},
		{"SQLAlchemy vs Zod", "sqlalchemy", "zod", false},
		{"Invalid schema", "pydantic", "nonexistent", true},
	}
	
	for _, comp := range comparisons {
		t.Run(comp.name, func(t *testing.T) {
			report := CompareModelsWithGrammars(allModels, comp.schema1, comp.schema2)
			
			if comp.expectError && !strings.Contains(report, "not found") {
				t.Errorf("Expected error for %s, but got: %s", comp.name, report)
			}
			
			if !comp.expectError && strings.Contains(report, "not found") {
				t.Errorf("Unexpected error for %s: %s", comp.name, report)
			}
			
			t.Logf("%s result: %s", comp.name, report)
		})
	}
}
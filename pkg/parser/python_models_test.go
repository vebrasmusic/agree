package parser

import "testing"

func TestParsePythonFiles(t *testing.T) {
	sql, pyd, err := ParsePythonFiles("../../test-data")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	// We now have more test data files, so expect more models
	if len(sql) < 1 {
		t.Fatalf("expected at least 1 sql model, got %d", len(sql))
	}
	if len(pyd) < 1 {
		t.Fatalf("expected at least 1 pydantic model, got %d", len(pyd))
	}
	
	// Check that we have the user model
	m, exists := sql["user"]
	if !exists {
		t.Fatalf("sql model 'user' not found. Available models: %v", getKeys(sql))
	}
	if _, ok := m.Fields["id"]; !ok {
		t.Fatalf("sql model missing field id")
	}
	
	pm, exists := pyd["user"]
	if !exists {
		t.Fatalf("pydantic model 'user' not found. Available models: %v", getKeys(pyd))
	}
	if _, ok := pm.Fields["username"]; !ok {
		t.Fatalf("pyd model missing field username")
	}
}

// Helper function to get map keys for debugging
func getKeys(m map[string]Model) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestCompareModels(t *testing.T) {
	sql, pyd, err := ParsePythonFiles("../../test-data")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	rep := CompareModels(sql, pyd)
	if rep == "" {
		t.Fatalf("expected report")
	}
}

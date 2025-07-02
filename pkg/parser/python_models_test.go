package parser

import "testing"

func TestParsePythonFiles(t *testing.T) {
	sql, pyd, err := ParsePythonFiles("../../test-data")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(sql) != 1 {
		t.Fatalf("expected 1 sql model, got %d", len(sql))
	}
	if len(pyd) != 1 {
		t.Fatalf("expected 1 pydantic model, got %d", len(pyd))
	}
	m := sql["user"]
	if _, ok := m.Fields["id"]; !ok {
		t.Fatalf("sql model missing field id")
	}
	pm := pyd["user"]
	if _, ok := pm.Fields["username"]; !ok {
		t.Fatalf("pyd model missing field username")
	}
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

package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLILegacyParsing tests the legacy parsing mode
func TestCLILegacyParsing(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "agree-test", "../../main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("agree-test")

	// Run the check command without grammar flag
	cmd = exec.Command("./agree-test", "check")
	cmd.Dir = "../.."
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Legacy parsing results") {
		t.Errorf("Expected legacy parsing output, got: %s", output)
	}
}

// TestCLIGrammarParsing tests the grammar-based parsing mode
func TestCLIGrammarParsing(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "agree-test", "../../main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("agree-test")

	// Run the check command with grammar flag
	cmd = exec.Command("./agree-test", "check", "--grammar")
	cmd.Dir = "../.."
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Grammar-based parsing results") {
		t.Errorf("Expected grammar parsing output, got: %s", output)
	}
	
	if !strings.Contains(output, "Available schema types") {
		t.Errorf("Expected schema types list, got: %s", output)
	}
}

// TestCLIGrammarFlag tests the grammar flag shorthand
func TestCLIGrammarFlag(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "agree-test", "../../main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("agree-test")

	// Run the check command with short grammar flag
	cmd = exec.Command("./agree-test", "check", "-g")
	cmd.Dir = "../.."
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Grammar-based parsing results") {
		t.Errorf("Expected grammar parsing output with -g flag, got: %s", output)
	}
}

// TestCLIErrorHandling tests error conditions
func TestCLIErrorHandling(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "agree-test", "../../main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("agree-test")

	// Test with invalid command
	cmd = exec.Command("./agree-test", "invalid-command")
	cmd.Dir = "../.."
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err == nil {
		t.Error("Expected error for invalid command, but got none")
	}
}

// TestCLIHelp tests the help command
func TestCLIHelp(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "agree-test", "../../main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("agree-test")

	// Test help command
	cmd = exec.Command("./agree-test", "--help")
	cmd.Dir = "../.."
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Help command failed: %v, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "agree") {
		t.Errorf("Expected help output to contain 'agree', got: %s", output)
	}
}

// TestGrammarFiles tests that grammar files exist and are valid
func TestGrammarFiles(t *testing.T) {
	grammarDir := "../../grammars"
	
	expectedGrammars := []string{"pydantic.json", "sqlalchemy.json", "zod.json"}
	
	for _, grammarFile := range expectedGrammars {
		grammarPath := filepath.Join(grammarDir, grammarFile)
		if _, err := os.Stat(grammarPath); os.IsNotExist(err) {
			t.Errorf("Grammar file %s does not exist", grammarPath)
		}
	}
}

// TestTestDataFiles tests that test data files exist
func TestTestDataFiles(t *testing.T) {
	testDataDir := "../../test-data"
	
	expectedFiles := []string{"pydantic.py", "sql_alchemy.py", "sqlalchemy_modern.py", "zod.ts"}
	
	for _, testFile := range expectedFiles {
		testPath := filepath.Join(testDataDir, testFile)
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Errorf("Test data file %s does not exist", testPath)
		}
	}
}

// TestGrammarLoadingIntegration tests that grammars load correctly
func TestGrammarLoadingIntegration(t *testing.T) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "agree-test", "../../main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("agree-test")

	// Run grammar parsing and check that all expected schema types are loaded
	cmd = exec.Command("./agree-test", "check", "--grammar")
	cmd.Dir = "../.."
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Grammar parsing failed: %v, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	
	// Check that all expected schema types are present
	expectedSchemaTypes := []string{"pydantic", "sqlalchemy", "zod"}
	for _, schemaType := range expectedSchemaTypes {
		if !strings.Contains(output, schemaType+":") {
			t.Errorf("Expected schema type %s in output, got: %s", schemaType, output)
		}
	}
}
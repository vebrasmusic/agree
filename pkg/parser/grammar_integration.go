package parser

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	ts "github.com/tree-sitter/go-tree-sitter"
	py "github.com/tree-sitter/tree-sitter-python/bindings/go"
	typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

// ParseFilesWithGrammars parses both Python and TypeScript files using the grammar engine
func ParseFilesWithGrammars(dir string, grammarDir string) (map[string]map[string]Model, error) {
	// Initialize grammar engine
	engine := NewGrammarEngine()

	// Load all grammar files from the grammar directory
	err := filepath.WalkDir(grammarDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		return engine.LoadGrammar(path)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load grammars: %w", err)
	}

	// Result map: schema_type -> nickname -> Model
	results := make(map[string]map[string]Model)

	// Initialize result maps for each loaded grammar
	for _, grammarName := range engine.GetGrammarNames() {
		results[grammarName] = make(map[string]Model)
	}

	// Parse files (Python support complete, TypeScript support planned)
	pythonLang := ts.NewLanguage(py.Language())
	// TODO: Add TypeScript support once tree-sitter binding imports are resolved
	typescriptLang := ts.NewLanguage(typescript.LanguageTypescript())

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".py" && ext != ".ts" && ext != ".tsx" {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Extract agree blocks
		blocks := extractAgreeBlocks(string(src))

		for _, block := range blocks {
			// Check if we have a grammar for this block type
			if models, exists := results[block.Type]; exists {
				// Handle both Python and TypeScript files
				var model Model
				var err error
				
				if ext == ".py" {
					model, err = engine.ParseModel([]byte(block.Code), block.Type, pythonLang)
				} else if ext == ".ts" || ext == ".tsx" {
					model, err = engine.ParseTypeScriptModel([]byte(block.Code), block.Type, typescriptLang)
				}
				
				if err != nil {
					return fmt.Errorf("%s: failed to parse %s block '%s': %w", path, block.Type, block.Nickname, err)
				}
				models[block.Nickname] = model
			}
		}

		return nil
	})

	return results, err
}

// CompareModelsWithGrammars compares models from different schema types
func CompareModelsWithGrammars(models map[string]map[string]Model, schemaType1, schemaType2 string) string {
	models1, ok1 := models[schemaType1]
	models2, ok2 := models[schemaType2]

	if !ok1 || !ok2 {
		return fmt.Sprintf("Schema types '%s' or '%s' not found", schemaType1, schemaType2)
	}

	return CompareModels(models1, models2)
}

// ParsePythonFilesWithGrammars is an alias for backward compatibility
func ParsePythonFilesWithGrammars(dir string, grammarDir string) (map[string]map[string]Model, error) {
	return ParseFilesWithGrammars(dir, grammarDir)
}

// Example usage function showing how to integrate with existing code
func ExampleGrammarUsage() error {
	// Parse files using grammars (supports Python and TypeScript)
	allModels, err := ParseFilesWithGrammars("test-data", "grammars")
	if err != nil {
		return err
	}

	// Compare SQLAlchemy vs Pydantic models
	report := CompareModelsWithGrammars(allModels, "sqlalchemy", "pydantic")
	fmt.Println("SQLAlchemy vs Pydantic comparison:")
	fmt.Println(report)

	// Compare Pydantic vs Zod models (cross-language)
	if zodReport := CompareModelsWithGrammars(allModels, "pydantic", "zod"); zodReport != "Schema types 'pydantic' or 'zod' not found" {
		fmt.Println("\nPydantic vs Zod comparison:")
		fmt.Println(zodReport)
	}

	return nil
}


/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vebrasmusic/agree/pkg/parser"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check your tracked schemas for missing changes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		useGrammar, _ := cmd.Flags().GetBool("grammar")
		
		if useGrammar {
			// Use new grammar-based parsing (supports both Python and TypeScript)
			allModels, err := parser.ParseFilesWithGrammars("test-data", "grammars")
			if err != nil {
				return err
			}
			
			// Compare sqlalchemy vs pydantic
			report := parser.CompareModelsWithGrammars(allModels, "sqlalchemy", "pydantic")
			fmt.Println("=== Grammar-based parsing results ===")
			fmt.Println("SQLAlchemy vs Pydantic:")
			fmt.Println(report)
			
			// Compare Pydantic vs Zod (cross-language)
			zodReport := parser.CompareModelsWithGrammars(allModels, "pydantic", "zod")
			fmt.Println("\nPydantic vs Zod (Python ↔ TypeScript):")
			fmt.Println(zodReport)
			
			// Show what grammars were loaded
			fmt.Println("\n=== Available schema types ===")
			for schemaType, models := range allModels {
				fmt.Printf("%s: %d models\n", schemaType, len(models))
			}
		} else {
			// Use original hardcoded parsing
			sqlModels, pydModels, err := parser.ParsePythonFiles("test-data")
			if err != nil {
				return err
			}
			report := parser.CompareModels(sqlModels, pydModels)
			fmt.Println("=== Legacy parsing results ===")
			fmt.Println(report)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Add flag to enable grammar-based parsing
	checkCmd.Flags().BoolP("grammar", "g", false, "Use grammar-based parsing instead of legacy hardcoded parsing")
}

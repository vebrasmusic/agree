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
		sqlModels, pydModels, err := parser.ParsePythonFiles("test-data")
		if err != nil {
			return err
		}
		report := parser.CompareModels(sqlModels, pydModels)
		fmt.Println(report)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

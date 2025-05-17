// Package main is the root command package for the personal-agent CLI.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "personal-agent",
	Short: "Personal Agent is a tool for managing personal documents",
	Long: `A CLI tool for managing and syncing personal documents
across different storage providers.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// Execute runs the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed
}

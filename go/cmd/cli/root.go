// Package main is the root command package for the personal-agent CLI.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "personal-agent",
	Short: "Personal Agent is a tool for managing personal documents",
	Long: `A CLI tool for managing and syncing personal documents
across different storage providers.`,
}

// appContext holds the application context including configuration
var appContext *AppContext

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
func Execute(ctx *AppContext) {
	appContext = ctx
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// GetAppContext returns the application context
func GetAppContext() *AppContext {
	return appContext
}

func init() {
	// Add global flags here if needed
}

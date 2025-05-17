// Package main is the entry point for the personal-agent CLI.
package main

import (
	"fmt"
	"os"

	"github.com/bonyuta0204/personal-agent/go/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Load and validate configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create application context with the loaded configuration
	appCtx := NewAppContext(cfg)

	// Execute the root command with the application context
	Execute(appCtx)
}

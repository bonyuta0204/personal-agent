// Package main implements the memory-related commands for the personal-agent CLI.
package main

import (
	"fmt"

	database "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/database"
	embeddingProvider "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/embedding"
	"github.com/bonyuta0204/personal-agent/go/internal/infrastructure/repository/postgres"
	storageFactory "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/storage"
	"github.com/bonyuta0204/personal-agent/go/internal/usecase/memory"
	"github.com/spf13/cobra"
)

// memoryCmd represents the memory command
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage memorys",
	Long:  `Commands for managing memorys in stores.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// This runs before any memory subcommand
		return nil
	},
}

var syncMemoryCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync memory",
	Long:  `Synchronize memorys`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repositories
		memoryRepo := postgres.NewMemoryRepository(db)

		// Initialize storage factory provider
		memoryStorageFactory := storageFactory.NewMemoryStorageFactory(ctx.Config.Memory.Repo)

		// Initialize embedding provider
		openaiProvider, err := embeddingProvider.NewOpenAIProvider()
		if err != nil {
			return fmt.Errorf("failed to create embedding provider: %w", err)
		}

		// Initialize sync use case
		syncUsecase := memory.NewSyncUsecase(memoryRepo, memoryStorageFactory, openaiProvider)

		// Execute the sync
		err = syncUsecase.Sync()
		if err != nil {
			return fmt.Errorf("sync failed: %w", err)
		}

		fmt.Println("Sync completed successfully")
		return nil
	},
}

var (
	// Flags for sync command
	dryRun bool
)

func init() {
	rootCmd.AddCommand(memoryCmd)
	memoryCmd.AddCommand(syncMemoryCmd)

	// Add flags for memory commands
	syncMemoryCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Perform a trial run with no changes made")
}

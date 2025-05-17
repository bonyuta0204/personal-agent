// Package main implements the document-related commands for the personal-agent CLI.
package main

import (
	"fmt"
	"strconv"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	database "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/database"
	"github.com/bonyuta0204/personal-agent/go/internal/infrastructure/repository/postgres"
	storageFactory "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/storage"
	"github.com/bonyuta0204/personal-agent/go/internal/usecase/document"
	"github.com/spf13/cobra"
)

// documentCmd represents the document command
var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents",
	Long:  `Commands for managing documents in stores.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// This runs before any document subcommand
		return nil
	},
}

// parseStoreID converts a string store ID to a uint
func parseStoreID(storeIDStr string) (model.StoreId, error) {
	id, err := strconv.ParseUint(storeIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid store ID format: %v", err)
	}
	return model.StoreId(id), nil
}

var syncDocumentCmd = &cobra.Command{
	Use:   "sync <store-id>",
	Short: "Sync documents from a store",
	Long:  `Synchronize documents from the specified store.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storeIDStr := args[0]
		storeID, err := parseStoreID(storeIDStr)
		if err != nil {
			return err
		}

		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repositories
		documentRepo := postgres.NewDocumentRepository(db)
		storeRepo := postgres.NewStoreRepository(db)

		// Initialize storage factory provider
		storageFactoryProvider := storageFactory.NewStorageFactoryProvider()

		// Initialize sync use case
		syncUsecase := document.NewSyncUsecase(storeRepo, documentRepo, storageFactoryProvider)

		// Execute the sync
		fmt.Printf("Starting sync for store ID: %d\n", storeID)

		err = syncUsecase.Sync(storeIDStr)
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
	rootCmd.AddCommand(documentCmd)
	documentCmd.AddCommand(syncDocumentCmd)

	// Add flags for document commands
	syncDocumentCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Perform a trial run with no changes made")
}

package main

import (
	"database/sql"
	"fmt"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	database "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/database"
	postgresRepo "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/repository/postgres"
	storeusecase "github.com/bonyuta0204/personal-agent/go/internal/usecase/store"
	"github.com/spf13/cobra"
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Manage document stores",
	Long:  `Commands for managing document storage locations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// This runs before any store subcommand
		return nil
	},
}

var createStoreCmd = &cobra.Command{
	Use:   "create [repository]",
	Short: "Create a new document store",
	Long:  `Create a new document store. For GitHub repositories, use the format "owner/repo".`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repository and use case
		repository := postgresRepo.NewStoreRepository(db)
		createUsecase := storeusecase.NewCreateUsecase(repository)

		// Create the store
		store, err := createUsecase.Create(repo)
		if err != nil {
			return fmt.Errorf("failed to create store: %w", err)
		}

		// Print success message
		fmt.Printf("Successfully created store with ID: %d\n", store.ID())
		fmt.Printf("Type: %s, Repository: %s\n", store.Type(), store.(*model.GitHubStore).Repo())

		return nil
	},
}

// listStoresResponse represents the structure of a store in the list
// This is a simplified version of the Store model for listing purposes
type listStoresResponse struct {
	ID   uint   `db:"id"`
	Type string `db:"type"`
	Repo string `db:"repo"`
}

var listStoresCmd = &cobra.Command{
	Use:   "list",
	Short: "List all document stores",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Query to list all stores
		query := `SELECT id, type, repo FROM stores ORDER BY id`

		var stores []listStoresResponse
		if err := db.Select(&stores, query); err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No document stores found")
				return nil
			}
			return fmt.Errorf("failed to list stores: %w", err)
		}

		if len(stores) == 0 {
			fmt.Println("No document stores found")
			return nil
		}

		// Print stores in a table format
		fmt.Println("ID  | Type    | Repository")
		fmt.Println("----|---------|-----------")
		for _, store := range stores {
			fmt.Printf("%-3d | %-7s | %s\n", store.ID, store.Type, store.Repo)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.AddCommand(createStoreCmd)
	storeCmd.AddCommand(listStoresCmd)

	// Add flags for store commands here if needed
	// Example:
	// createStoreCmd.Flags().StringP("type", "t", "github", "Type of the store (e.g., github)")
}

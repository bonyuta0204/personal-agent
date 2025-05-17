package main

import (
	"fmt"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
	storeusecase "github.com/bonyuta0204/personal-agent/go/internal/usecase/store"
	"github.com/spf13/cobra"
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Manage document stores",
	Long:  `Commands for managing document storage locations.`,
}

var createStoreCmd = &cobra.Command{
	Use:   "create [repository]",
	Short: "Create a new document store",
	Long:  `Create a new document store. For GitHub repositories, use the format "owner/repo".`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]
		ctx := GetAppContext()

		// Log database connection info (for debugging, remove in production)
		cfg := ctx.Config
		fmt.Printf("Connecting to database: %s@%s:%s/%s\n", 
			cfg.Database.User, 
			cfg.Database.Host, 
			cfg.Database.Port, 
			cfg.Database.Name)

		// Initialize dependencies with configuration
		// In a real implementation, you would use the config to initialize your database connection
		// For example:
		// db, err := sql.Open("postgres", cfg.Database.GetDSN())
		// if err != nil {
		//     return fmt.Errorf("failed to connect to database: %v", err)
		// }
		// defer db.Close()
		
		// Initialize repository and use case with database connection
		_ = repository.StoreRepository(nil) // Pass db here
		_ = storeusecase.NewCreateUsecase(nil)


		fmt.Printf("Creating store for repository: %s\n", repo)
		// store, err := createUsecase.Create(repo)
		// if err != nil {
		//     return fmt.Errorf("failed to create store: %v", err)
		// }
		// fmt.Printf("Successfully created store with ID: %d\n", store.ID())
		
		return nil
	},
}

var listStoresCmd = &cobra.Command{
	Use:   "list",
	Short: "List all document stores",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := GetAppContext()
		
		// In a real implementation, you would use the config to connect to the database
		// and list all stores
		fmt.Printf("Listing all document stores from database: %s\n", 
			ctx.Config.Database.Name)
		
		// Example of accessing database configuration
		// db, err := sql.Open("postgres", ctx.Config.Database.GetDSN())
		// if err != nil {
		//     return fmt.Errorf("failed to connect to database: %v", err)
		// }
		// defer db.Close()
		
		// Implementation would go here
		// stores, err := storeRepository.List()
		// if err != nil {
		//     return fmt.Errorf("failed to list stores: %v", err)
		// }
		
		// for _, store := range stores {
		//     fmt.Printf("- %s (ID: %d)\n", store.Name(), store.ID())
		// }
		
		fmt.Println("No stores found (not implemented yet)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.AddCommand(createStoreCmd)
	storeCmd.AddCommand(listStoresCmd)

	// Add flags for store commands here if needed
}

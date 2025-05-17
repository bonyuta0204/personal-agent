// Package main implements the document-related commands for the personal-agent CLI.
package main

import (
	"fmt"

	"github.com/bonyuta0204/personal-agent/go/internal/usecase/document"
	"github.com/spf13/cobra"
)

// documentCmd represents the document command
var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents",
	Long:  `Commands for managing documents in stores.`,
}

var syncDocumentCmd = &cobra.Command{
	Use:   "sync [store-id]",
	Short: "Sync documents from a store",
	Long:  `Synchronize documents from the specified store.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		storeID := args[0]

		// Initialize dependencies (in a real app, this would be handled by a dependency injection container)
		// For now, we'll just show the command structure
		_ = document.NewSyncUsecase(nil, nil, nil)

		// In a real implementation, we would call the use case here:
		// err := syncUsecase.Sync(storeID)
		// if err != nil {
		//     return fmt.Errorf("failed to sync documents: %v", err)
		// }

		fmt.Printf("Would sync documents from store ID: %s\n", storeID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(documentCmd)
	documentCmd.AddCommand(syncDocumentCmd)

	// Add flags for document commands here if needed
}

package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	database "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/database"
	"github.com/bonyuta0204/personal-agent/go/internal/infrastructure/embedding/openai"
	postgresRepo "github.com/bonyuta0204/personal-agent/go/internal/infrastructure/repository/postgres"
	memoryusecase "github.com/bonyuta0204/personal-agent/go/internal/usecase/memory"
	"github.com/spf13/cobra"
)

// memoryCmd represents the memory command
var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage AI agent memories",
	Long:  `Commands for managing AI agent memory storage.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var saveMemoryCmd = &cobra.Command{
	Use:   "save [path] [content]",
	Short: "Save a memory",
	Long:  `Save content as a memory with optional tags.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		content := args[1]
		
		// Get tags from flag
		tags, _ := cmd.Flags().GetStringSlice("tags")
		
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repository and embedding provider
		memoryRepo := postgresRepo.NewMemoryRepository(db)
		embeddingProvider := openai.NewProvider(ctx.Config.OpenAI.APIKey)
		
		// Initialize use case
		saveUsecase := memoryusecase.NewSaveUseCase(memoryRepo, embeddingProvider)

		// Save the memory
		memory, err := saveUsecase.Execute(cmd.Context(), &memoryusecase.SaveMemoryRequest{
			Path:    path,
			Content: content,
			Tags:    tags,
		})
		if err != nil {
			return fmt.Errorf("failed to save memory: %w", err)
		}

		fmt.Printf("Successfully saved memory with ID: %s\n", memory.ID)
		fmt.Printf("Path: %s\n", memory.Path)
		if len(memory.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(memory.Tags, ", "))
		}

		return nil
	},
}

var listMemoriesCmd = &cobra.Command{
	Use:   "list",
	Short: "List memories",
	Long:  `List stored memories with pagination.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get pagination flags
		offset, _ := cmd.Flags().GetInt("offset")
		limit, _ := cmd.Flags().GetInt("limit")
		
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repository and use case
		memoryRepo := postgresRepo.NewMemoryRepository(db)
		retrieveUsecase := memoryusecase.NewRetrieveUseCase(memoryRepo)

		// List memories
		memories, err := retrieveUsecase.List(cmd.Context(), offset, limit)
		if err != nil {
			return fmt.Errorf("failed to list memories: %w", err)
		}

		if len(memories) == 0 {
			fmt.Println("No memories found")
			return nil
		}

		// Print memories
		fmt.Println("ID   | Path                           | Tags")
		fmt.Println("-----|--------------------------------|-----")
		for _, memory := range memories {
			tagsStr := strings.Join(memory.Tags, ", ")
			if len(tagsStr) > 30 {
				tagsStr = tagsStr[:27] + "..."
			}
			pathStr := memory.Path
			if len(pathStr) > 30 {
				pathStr = pathStr[:27] + "..."
			}
			fmt.Printf("%-4s | %-30s | %s\n", memory.ID, pathStr, tagsStr)
		}

		return nil
	},
}

var getMemoryCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a memory by ID",
	Long:  `Retrieve and display a specific memory by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := model.MemoryId(args[0])
		
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repository and use case
		memoryRepo := postgresRepo.NewMemoryRepository(db)
		retrieveUsecase := memoryusecase.NewRetrieveUseCase(memoryRepo)

		// Get memory
		memory, err := retrieveUsecase.GetByID(cmd.Context(), id)
		if err != nil {
			return fmt.Errorf("failed to get memory: %w", err)
		}

		// Print memory details
		fmt.Printf("ID: %s\n", memory.ID)
		fmt.Printf("Path: %s\n", memory.Path)
		fmt.Printf("Content: %s\n", memory.Content)
		if len(memory.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(memory.Tags, ", "))
		}
		fmt.Printf("SHA: %s\n", memory.SHA)
		fmt.Printf("Modified At: %s\n", memory.ModifiedAt.Format(time.RFC3339))
		fmt.Printf("Created At: %s\n", memory.CreatedAt.Format(time.RFC3339))

		return nil
	},
}

var searchMemoriesCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search memories by text",
	Long:  `Search memories using semantic similarity.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		
		// Get limit flag
		limit, _ := cmd.Flags().GetInt("limit")
		
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repository and embedding provider
		memoryRepo := postgresRepo.NewMemoryRepository(db)
		embeddingProvider := openai.NewProvider(ctx.Config.OpenAI.APIKey)
		
		// Initialize use case
		searchUsecase := memoryusecase.NewSearchUseCase(memoryRepo, embeddingProvider)

		// Search memories
		memories, err := searchUsecase.SearchByText(cmd.Context(), &memoryusecase.SearchByTextRequest{
			Query: query,
			Limit: limit,
		})
		if err != nil {
			return fmt.Errorf("failed to search memories: %w", err)
		}

		if len(memories) == 0 {
			fmt.Printf("No memories found for query: %s\n", query)
			return nil
		}

		// Print search results
		fmt.Printf("Found %d memories for query: %s\n\n", len(memories), query)
		for i, memory := range memories {
			fmt.Printf("%d. ID: %s\n", i+1, memory.ID)
			fmt.Printf("   Path: %s\n", memory.Path)
			if len(memory.Tags) > 0 {
				fmt.Printf("   Tags: %s\n", strings.Join(memory.Tags, ", "))
			}
			fmt.Printf("   Content: %.100s", memory.Content)
			if len(memory.Content) > 100 {
				fmt.Printf("...")
			}
			fmt.Printf("\n\n")
		}

		return nil
	},
}

var searchMemoriesByTagsCmd = &cobra.Command{
	Use:   "search-tags [tags...]",
	Short: "Search memories by tags",
	Long:  `Search memories that contain all specified tags.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tags := args
		
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repository and use case
		memoryRepo := postgresRepo.NewMemoryRepository(db)
		searchUsecase := memoryusecase.NewSearchUseCase(memoryRepo, nil)

		// Search memories by tags
		memories, err := searchUsecase.SearchByTags(cmd.Context(), &memoryusecase.SearchByTagsRequest{
			Tags: tags,
		})
		if err != nil {
			return fmt.Errorf("failed to search memories by tags: %w", err)
		}

		if len(memories) == 0 {
			fmt.Printf("No memories found with tags: %s\n", strings.Join(tags, ", "))
			return nil
		}

		// Print search results
		fmt.Printf("Found %d memories with tags: %s\n\n", len(memories), strings.Join(tags, ", "))
		for i, memory := range memories {
			fmt.Printf("%d. ID: %s\n", i+1, memory.ID)
			fmt.Printf("   Path: %s\n", memory.Path)
			fmt.Printf("   Tags: %s\n", strings.Join(memory.Tags, ", "))
			fmt.Printf("   Content: %.100s", memory.Content)
			if len(memory.Content) > 100 {
				fmt.Printf("...")
			}
			fmt.Printf("\n\n")
		}

		return nil
	},
}

var deleteMemoryCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a memory",
	Long:  `Delete a memory by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := model.MemoryId(args[0])
		
		ctx := GetAppContext()

		// Initialize database connection
		db, err := database.NewDBConnection(&ctx.Config.Database)
		if err != nil {
			return fmt.Errorf("database connection error: %w", err)
		}
		defer database.CloseDB(db)

		// Initialize repository and use case
		memoryRepo := postgresRepo.NewMemoryRepository(db)
		deleteUsecase := memoryusecase.NewDeleteUseCase(memoryRepo)

		// Delete memory
		err = deleteUsecase.Execute(cmd.Context(), id)
		if err != nil {
			return fmt.Errorf("failed to delete memory: %w", err)
		}

		fmt.Printf("Successfully deleted memory with ID: %s\n", id)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(memoryCmd)
	memoryCmd.AddCommand(saveMemoryCmd)
	memoryCmd.AddCommand(listMemoriesCmd)
	memoryCmd.AddCommand(getMemoryCmd)
	memoryCmd.AddCommand(searchMemoriesCmd)
	memoryCmd.AddCommand(searchMemoriesByTagsCmd)
	memoryCmd.AddCommand(deleteMemoryCmd)

	// Add flags
	saveMemoryCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags for the memory")
	
	listMemoriesCmd.Flags().IntP("offset", "o", 0, "Offset for pagination")
	listMemoriesCmd.Flags().IntP("limit", "l", 10, "Limit for pagination")
	
	searchMemoriesCmd.Flags().IntP("limit", "l", 10, "Limit for search results")
}
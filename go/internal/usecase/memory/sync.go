package memory

import (
	"fmt"
	"log"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/embedding"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/storage"
)

type SyncUsecase struct {
	memoryRepo           repository.MemoryRepository
	memoryStorageFactory storage.MemoryStorageFactory
	embeddingProvider    embedding.EmbeddingProvider
}

// NewSyncUsecase creates a new SyncUsecase instance
func NewSyncUsecase(memoryRepo repository.MemoryRepository, memoryStorageFactory storage.MemoryStorageFactory, embeddingProvider embedding.EmbeddingProvider) *SyncUsecase {
	return &SyncUsecase{
		memoryRepo:           memoryRepo,
		memoryStorageFactory: memoryStorageFactory,
		embeddingProvider:    embeddingProvider,
	}
}

func (u *SyncUsecase) Sync(storeId string) error {
	// Convert string storeId to model.StoreId (uint)
	var id model.StoreId
	_, err := fmt.Sscanf(storeId, "%d", &id)
	if err != nil {
		return fmt.Errorf("invalid store ID format: %v", err)
	}

	// Get the storage
	storage, err := u.memoryStorageFactory.CreateMemoryStorage()
	if err != nil {
		return fmt.Errorf("failed to get storage factory: %w", err)
	}

	// Get all document entries from storage
	entries, err := storage.GetMemoryEntries()
	if err != nil {
		return fmt.Errorf("failed to get memory entries: %w", err)
	}

	// Fetch all documents from storage
	var memories []*model.Memory

	for _, entry := range entries {
		memory, err := storage.FetchMemory(entry.Path)
		if err != nil {
			log.Printf("failed to fetch memory %s: %v", entry.Path, err)
			continue
		}
		if memory != nil {
			memories = append(memories, memory)
		}
	}

	// Find unchanged memories
	existingSHAs, err := u.memoryRepo.FindExistingSHAs(memories)
	if err != nil {
		return fmt.Errorf("failed to find unchanged memories: %w", err)
	}

	// Create a map for quick lookup of unchanged document IDs
	existingSHASet := make(map[string]bool, len(existingSHAs))
	for _, sha := range existingSHAs {
		existingSHASet[sha] = true
	}

	// Save only changed memories
	savedCount := 0
	for _, mem := range memories {
		if mem == nil {
			continue
		}

		if !existingSHASet[mem.SHA] {
			// create embedding
			embedding, err := u.embeddingProvider.Embed(mem.Content)
			if err != nil {
				log.Printf("failed to create embedding for memory %s: %v", mem.Path, err)
				continue
			}
			mem.Embedding = embedding
			if err := u.memoryRepo.SaveMemory(mem); err != nil {
				log.Printf("failed to save memory %s: %v", mem.Path, err)
				continue
			}
			savedCount++
		} else {
		}
	}

	log.Printf("sync completed: %d memories processed, %d memories saved", len(memories), savedCount)

	return nil
}

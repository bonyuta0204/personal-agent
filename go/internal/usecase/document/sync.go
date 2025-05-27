package document

import (
	"fmt"
	"log"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/embedding"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/storage"
)

type SyncUsecase struct {
	storeRepo              repository.StoreRepository
	documentRepo           repository.DocumentRepository
	storageFactoryProvider storage.StorageFactoryProvider
	embeddingProvider      embedding.EmbeddingProvider
}

// NewSyncUsecase creates a new SyncUsecase instance
func NewSyncUsecase(storeRepo repository.StoreRepository, documentRepo repository.DocumentRepository, factoryProvider storage.StorageFactoryProvider, embeddingProvider embedding.EmbeddingProvider) *SyncUsecase {
	return &SyncUsecase{
		storeRepo:              storeRepo,
		documentRepo:           documentRepo,
		storageFactoryProvider: factoryProvider,
		embeddingProvider:      embeddingProvider,
	}
}

func (u *SyncUsecase) Sync(storeId string) error {
	// Convert string storeId to model.StoreId (uint)
	var id model.StoreId
	_, err := fmt.Sscanf(storeId, "%d", &id)
	if err != nil {
		return fmt.Errorf("invalid store ID format: %v", err)
	}

	store, err := u.storeRepo.GetStore(id)
	if err != nil {
		return err
	}

	// Get the appropriate storage factory for this store type
	factory, err := u.storageFactoryProvider.GetFactory(store.Type())
	if err != nil {
		return fmt.Errorf("failed to get storage factory: %w", err)
	}

	// Create the storage instance
	storage, err := factory.CreateStorage(store)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// Get all document entries from storage
	entries, err := storage.GetDocumentEntries()
	if err != nil {
		return fmt.Errorf("failed to get document entries: %w", err)
	}

	// Fetch all documents from storage
	var documents []*model.Document

	for _, entry := range entries {
		document, err := storage.FetchDocument(store.ID(), entry.Path)
		if err != nil {
			log.Printf("failed to fetch document %s: %v", entry.Path, err)
			continue
		}
		if document != nil {
			// set document tags from content
			document.SetTagsFromContent()
			documents = append(documents, document)
		}
	}

	// Find unchanged documents
	existingSHAs, err := u.documentRepo.FindExistingSHAs(documents)
	if err != nil {
		return fmt.Errorf("failed to find unchanged documents: %w", err)
	}

	// Create a map for quick lookup of unchanged document IDs
	existingSHASet := make(map[string]bool, len(existingSHAs))
	for _, sha := range existingSHAs {
		existingSHASet[sha] = true
	}

	// Save only changed documents
	savedCount := 0
	for _, doc := range documents {
		if doc == nil {
			continue
		}

		if !existingSHASet[doc.SHA] {
			// create embedding
			embedding, err := u.embeddingProvider.Embed(doc.Content)
			if err != nil {
				log.Printf("failed to create embedding for document %s: %v", doc.Path, err)
				continue
			}
			doc.Embedding = embedding
			if err := u.documentRepo.SaveDocument(doc); err != nil {
				log.Printf("failed to save document %s: %v", doc.Path, err)
				continue
			}
			savedCount++
		} else {
		}
	}

	log.Printf("sync completed: %d documents processed, %d documents saved", len(documents), savedCount)

	return nil
}

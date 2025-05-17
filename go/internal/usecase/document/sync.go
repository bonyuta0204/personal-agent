package document

import (
	"fmt"
	"log"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/storage"
)

type SyncUsecase struct {
	storeRepo              repository.StoreRepository
	documentRepo           repository.DocumentRepository
	storageFactoryProvider storage.StorageFactoryProvider
}

// NewSyncUsecase creates a new SyncUsecase instance
func NewSyncUsecase(storeRepo repository.StoreRepository, documentRepo repository.DocumentRepository, factoryProvider storage.StorageFactoryProvider) *SyncUsecase {
	return &SyncUsecase{
		storeRepo:              storeRepo,
		documentRepo:           documentRepo,
		storageFactoryProvider: factoryProvider,
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

	entries, err := storage.GetDocumentEntries()
	if err != nil {
		return fmt.Errorf("failed to get all paths: %w", err)
	}

	for _, entry := range entries {
		document, err := storage.FetchDocument(store.ID(), entry.Path)
		if err != nil {
			log.Printf("failed to fetch document: %s", err)
		}

		if err := u.documentRepo.SaveDocument(document); err != nil {
			log.Printf("failed to save document: %s", err)
		}

	}

	return nil
}

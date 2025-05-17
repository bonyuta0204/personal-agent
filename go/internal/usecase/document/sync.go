package document

import (
	"fmt"

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
	store, err := u.storeRepo.GetStore(storeId)
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

	paths, err := storage.GetAllPaths()
	if err != nil {
		return fmt.Errorf("failed to get all paths: %w", err)
	}

	for _, path := range paths {
		document, err := storage.FetchDocument(path)
		if err != nil {
			return fmt.Errorf("failed to fetch document: %w", err)
		}

		if err := u.documentRepo.SaveDocument(document); err != nil {
			return fmt.Errorf("failed to save document: %w", err)
		}

	}

	return nil
}

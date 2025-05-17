package storage

import (
	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
)

type Storage interface {
	SaveDocument(document *model.Document) error
	SaveMemory(memory *model.Memory) error
	FetchDocument(storeId model.StoreId, path string) (*model.Document, error)
	FetchMemory(path string) (*model.Memory, error)
	GetDocumentEntries() ([]model.DocumentEntry, error)
}

type StorageFactory interface {
	CreateStorage(store model.DocumentStore) (Storage, error)
}
type StorageFactoryProvider interface {
	GetFactory(storeType string) (StorageFactory, error)
}

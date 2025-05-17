package storage

import (
	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
)

type Storage interface {
	SaveDocument(document *model.Document) error
	SaveMemory(memory *model.Memory) error
	FetchDocument(path string) (*model.Document, error)
	FetchMemory(path string) (*model.Memory, error)
	GetAllPaths() ([]string, error)
}

type StorageFactory interface {
	CreateStorage(store model.DocumentStore) (Storage, error)
}
type StorageFactoryProvider interface {
	GetFactory(storeType string) (StorageFactory, error)
}

package storage

import (
	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	port "github.com/bonyuta0204/personal-agent/go/internal/domain/port/storage"
)

// StorageFactoryProvider provides the appropriate storage factory based on store type
type StorageFactoryProvider struct {
	githubFactory *GitHubStorageFactory
}

// NewStorageFactoryProvider creates a new storage factory provider
func NewStorageFactoryProvider() *StorageFactoryProvider {
	return &StorageFactoryProvider{
		githubFactory: NewGitHubStorageFactory(),
	}
}

// GetFactory returns the appropriate storage factory for the given store type
func (p *StorageFactoryProvider) GetFactory(storeType string) (port.StorageFactory, error) {
	switch storeType {
	case model.StoreTypeGitHub:
		return p.githubFactory, nil
	default:
		return nil, model.ErrUnsupportedStoreType
	}
}

// Ensure MemoryStorageFactory implements port.MemoryStorageFactory
var _ port.MemoryStorageFactory = (*MemoryStorageFactory)(nil)

// MemoryStorageFactory implements the MemoryStorageFactory interface for memory
type MemoryStorageFactory struct {
	repo string
}

// NewMemoryStorageFactory creates a new memory storage factory
func NewMemoryStorageFactory(repo string) *MemoryStorageFactory {
	return &MemoryStorageFactory{
		repo: repo,
	}
}

// CreateMemoryStorage creates a new memory storage instance
func (f *MemoryStorageFactory) CreateMemoryStorage() (port.Storage, error) {
	return NewGitHubStorage(f.repo)
}

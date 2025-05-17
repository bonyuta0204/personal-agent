package storage

import (
	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/go/internal/usecase/document"
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
func (p *StorageFactoryProvider) GetFactory(storeType string) (document.StorageFactory, error) {
	switch storeType {
	case model.StoreTypeGitHub:
		return p.githubFactory, nil
	default:
		return nil, model.ErrUnsupportedStoreType
	}
}

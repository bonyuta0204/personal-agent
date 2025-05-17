package storage

import (
	"fmt"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	port "github.com/bonyuta0204/personal-agent/go/internal/domain/port/storage"
)

// GitHubStorageFactory implements the StorageFactory interface for GitHub
type GitHubStorageFactory struct{}

// NewGitHubStorageFactory creates a new GitHub storage factory
func NewGitHubStorageFactory() *GitHubStorageFactory {
	return &GitHubStorageFactory{}
}

// CreateStorage creates a new GitHub storage instance
func (f *GitHubStorageFactory) CreateStorage(store model.DocumentStore) (port.Storage, error) {
	if store.Type() != model.StoreTypeGitHub {
		return nil, fmt.Errorf("unsupported store type: %s", store.Type())
	}

	// Type assert to GitHubStore to access GitHub-specific fields
	githubStore, ok := store.(*model.GitHubStore)
	if !ok {
		return nil, fmt.Errorf("invalid store type for GitHub")
	}

	return NewGitHubStorage(githubStore.Repo())
}

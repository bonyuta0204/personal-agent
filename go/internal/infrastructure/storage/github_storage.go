package storage

import (
	"fmt"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/storage"
)

// GitHubStorage implements the storage.Storage interface for GitHub
// This is an infrastructure layer component
type GitHubStorage struct {
	// Add any GitHub client or configuration here
	repo string // owner/repo format
}

// NewGitHubStorage creates a new GitHub storage instance
func NewGitHubStorage(repo string) *GitHubStorage {
	return &GitHubStorage{
		repo: repo,
	}
}

// SaveDocument implements the Storage interface
func (s *GitHubStorage) SaveDocument(document *model.Document) error {
	// Implement GitHub-specific document saving logic here
	// This is where you would interact with the GitHub API
	return nil
}

// SaveMemory implements the Storage interface
func (s *GitHubStorage) SaveMemory(memory *model.Memory) error {
	// Implement GitHub-specific memory saving logic here
	return nil
}

// FetchDocument implements the Storage interface
func (s *GitHubStorage) FetchDocument(path string) (*model.Document, error) {
	// Implement GitHub-specific document fetching logic here
	return nil, nil
}

// FetchMemory implements the Storage interface
func (s *GitHubStorage) FetchMemory(path string) (*model.Memory, error) {
	// Implement GitHub-specific memory fetching logic here
	return nil, nil
}

// GitHubStorageFactory implements the StorageFactory interface for GitHub
type GitHubStorageFactory struct {
	// Add any dependencies needed for GitHub storage (like API clients)
}

// NewGitHubStorageFactory creates a new GitHub storage factory
func NewGitHubStorageFactory() *GitHubStorageFactory {
	return &GitHubStorageFactory{}
}

// CreateStorage creates a new GitHub storage instance
func (f *GitHubStorageFactory) CreateStorage(store model.DocumentStore) (storage.Storage, error) {
	if store.Type() != model.StoreTypeGitHub {
		return nil, fmt.Errorf("unsupported store type: %s", store.Type())
	}

	// Type assert to GitHubStore to access GitHub-specific fields
	githubStore, ok := store.(*model.GitHubStore)
	if !ok {
		return nil, fmt.Errorf("invalid store type for GitHub")
	}

	return NewGitHubStorage(githubStore.Repo()), nil
}

package model

import "errors"

type StoreId uint

const (
	StoreTypeGitHub = "github"
)

var (
	// ErrUnsupportedStoreType is returned when an unsupported store type is provided
	ErrUnsupportedStoreType = errors.New("unsupported store type")
)

type DocumentStore interface {
	ID() StoreId
	Type() string
}

type GitHubStore struct {
	id   StoreId
	repo string // owner/repo
}

// NewGitHubStore creates a new GitHub store instance
func NewGitHubStore(id StoreId, repo string) *GitHubStore {
	return &GitHubStore{
		id:   id,
		repo: repo,
	}
}

// ID returns the store ID
func (s *GitHubStore) ID() StoreId {
	return s.id
}

// Type returns the store type
func (s *GitHubStore) Type() string {
	return StoreTypeGitHub
}

// Repo returns the GitHub repository in owner/repo format
func (s *GitHubStore) Repo() string {
	return s.repo
}

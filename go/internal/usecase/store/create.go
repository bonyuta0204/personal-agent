package store

import (
	"errors"
	"fmt"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
)

type CreateUsecase struct {
	storeRepo repository.StoreRepository
}

func NewCreateUsecase(storeRepo repository.StoreRepository) *CreateUsecase {
	return &CreateUsecase{
		storeRepo: storeRepo,
	}
}

// Create creates a new document store
// Currently only GitHub repositories are supported in the format "owner/repo"
// Returns the created store with its generated ID
func (u *CreateUsecase) Create(repo string) (model.DocumentStore, error) {
	if repo == "" {
		return nil, errors.New("repository cannot be empty")
	}

	// Create a new store with a temporary ID (0 for auto-increment)
	// The actual ID will be assigned by the database
	tempStore := model.NewGitHubStore(0, repo)
	
	// Save the store and get the version with the generated ID
	createdStore, err := u.storeRepo.CreateStore(tempStore)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}
	
	return createdStore, nil
}

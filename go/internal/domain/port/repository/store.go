package repository

import (
	"errors"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
)

// ErrStoreNotFound is returned when a store is not found in the repository
var ErrStoreNotFound = errors.New("store not found")

type StoreRepository interface {
	GetStore(storeId model.StoreId) (model.DocumentStore, error)
	CreateStore(store model.DocumentStore) (model.DocumentStore, error)
}

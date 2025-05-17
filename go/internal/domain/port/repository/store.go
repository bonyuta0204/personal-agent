package repository

import "github.com/bonyuta0204/personal-agent/go/internal/domain/model"

type StoreRepository interface {
	GetStore(storeId model.StoreId) (model.DocumentStore, error)
	CreateStore(store model.DocumentStore) (model.DocumentStore, error)
}

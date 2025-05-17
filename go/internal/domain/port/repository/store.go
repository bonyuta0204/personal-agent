package repository

import "github.com/bonyuta0204/personal-agent/go/internal/domain/model"

type StoreRepository interface {
	GetStore(storeId string) (model.DocumentStore, error)
}

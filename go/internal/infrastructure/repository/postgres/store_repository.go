package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	repo "github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
	"github.com/jmoiron/sqlx"
)

// Ensure storeRepository implements repo.StoreRepository
var _ repo.StoreRepository = (*storeRepository)(nil)

type storeRepository struct {
	db *sqlx.DB
}

// NewStoreRepository creates a new PostgreSQL store repository
func NewStoreRepository(db *sqlx.DB) repo.StoreRepository {
	return &storeRepository{db: db}
}

// GetStore retrieves a store by ID
func (r *storeRepository) GetStore(storeID model.StoreId) (model.DocumentStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var store struct {
		ID   uint   `db:"id"`
		Type string `db:"type"`
		Repo string `db:"repo"`
	}

	query := `SELECT id, type, repo FROM stores WHERE id = $1`
	err := r.db.GetContext(ctx, &store, query, storeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrStoreNotFound
		}
		return nil, err
	}

	switch store.Type {
	case model.StoreTypeGitHub:
		return model.NewGitHubStore(model.StoreId(store.ID), store.Repo), nil
	default:
		return nil, model.ErrUnsupportedStoreType
	}
}

// CreateStore creates a new store
func (r *storeRepository) CreateStore(store model.DocumentStore) (model.DocumentStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch s := store.(type) {
	case *model.GitHubStore:
		var id uint
		query := `INSERT INTO stores (type, repo) VALUES ($1, $2) RETURNING id`
		err := r.db.QueryRowContext(ctx, query, store.Type(), s.Repo()).Scan(&id)
		if err != nil {
			return nil, err
		}
		return model.NewGitHubStore(model.StoreId(id), s.Repo()), nil
	default:
		return nil, model.ErrUnsupportedStoreType
	}
}

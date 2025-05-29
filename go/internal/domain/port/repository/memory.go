package repository

import (
	"context"

	"github.com/bonyuta0204/personal-agent/internal/domain/model"
)

type MemoryRepository interface {
	Create(ctx context.Context, memory *model.Memory) (*model.Memory, error)
	GetByID(ctx context.Context, id model.MemoryId) (*model.Memory, error)
	GetByPath(ctx context.Context, path string) (*model.Memory, error)
	List(ctx context.Context, offset, limit int) ([]*model.Memory, error)
	Update(ctx context.Context, memory *model.Memory) error
	Delete(ctx context.Context, id model.MemoryId) error
	SearchByEmbedding(ctx context.Context, embedding []float64, limit int) ([]*model.Memory, error)
	SearchByTags(ctx context.Context, tags []string) ([]*model.Memory, error)
}
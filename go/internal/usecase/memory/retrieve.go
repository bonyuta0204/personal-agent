package memory

import (
	"context"
	"fmt"

	"github.com/bonyuta0204/personal-agent/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/internal/domain/port/repository"
)

type RetrieveUseCase struct {
	memoryRepo repository.MemoryRepository
}

func NewRetrieveUseCase(memoryRepo repository.MemoryRepository) *RetrieveUseCase {
	return &RetrieveUseCase{
		memoryRepo: memoryRepo,
	}
}

func (uc *RetrieveUseCase) GetByID(ctx context.Context, id model.MemoryId) (*model.Memory, error) {
	memory, err := uc.memoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory by ID: %w", err)
	}
	return memory, nil
}

func (uc *RetrieveUseCase) GetByPath(ctx context.Context, path string) (*model.Memory, error) {
	memory, err := uc.memoryRepo.GetByPath(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory by path: %w", err)
	}
	return memory, nil
}

func (uc *RetrieveUseCase) List(ctx context.Context, offset, limit int) ([]*model.Memory, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	memories, err := uc.memoryRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list memories: %w", err)
	}
	return memories, nil
}
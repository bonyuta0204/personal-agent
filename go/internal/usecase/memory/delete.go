package memory

import (
	"context"
	"fmt"

	"github.com/bonyuta0204/personal-agent/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/internal/domain/port/repository"
)

type DeleteUseCase struct {
	memoryRepo repository.MemoryRepository
}

func NewDeleteUseCase(memoryRepo repository.MemoryRepository) *DeleteUseCase {
	return &DeleteUseCase{
		memoryRepo: memoryRepo,
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id model.MemoryId) error {
	err := uc.memoryRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}
	return nil
}
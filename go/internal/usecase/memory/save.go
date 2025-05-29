package memory

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/bonyuta0204/personal-agent/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/internal/domain/port/embedding"
	"github.com/bonyuta0204/personal-agent/internal/domain/port/repository"
	"github.com/bonyuta0204/personal-agent/internal/infrastructure/util"
)

type SaveUseCase struct {
	memoryRepo      repository.MemoryRepository
	embeddingClient embedding.Provider
}

func NewSaveUseCase(memoryRepo repository.MemoryRepository, embeddingClient embedding.Provider) *SaveUseCase {
	return &SaveUseCase{
		memoryRepo:      memoryRepo,
		embeddingClient: embeddingClient,
	}
}

type SaveMemoryRequest struct {
	Path       string
	Content    string
	Tags       []string
	ModifiedAt *time.Time
}

func (uc *SaveUseCase) Execute(ctx context.Context, req *SaveMemoryRequest) (*model.Memory, error) {
	// Calculate SHA for content
	sha := util.CalculateSHA(req.Content)

	// Check if memory already exists with same path
	existingMemory, err := uc.memoryRepo.GetByPath(ctx, req.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing memory: %w", err)
	}

	// If memory exists and SHA is the same, return existing memory
	if existingMemory != nil && existingMemory.SHA == sha {
		return existingMemory, nil
	}

	// Generate embedding for the content
	embeddingVector, err := uc.embeddingClient.GenerateEmbedding(ctx, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	modifiedAt := time.Now()
	if req.ModifiedAt != nil {
		modifiedAt = *req.ModifiedAt
	}

	memory := &model.Memory{
		Path:       req.Path,
		Content:    req.Content,
		Embedding:  embeddingVector,
		Tags:       req.Tags,
		SHA:        sha,
		ModifiedAt: modifiedAt,
	}

	// Update existing memory or create new one
	if existingMemory != nil {
		memory.ID = existingMemory.ID
		memory.CreatedAt = existingMemory.CreatedAt
		err = uc.memoryRepo.Update(ctx, memory)
		if err != nil {
			return nil, fmt.Errorf("failed to update memory: %w", err)
		}
	} else {
		memory, err = uc.memoryRepo.Create(ctx, memory)
		if err != nil {
			return nil, fmt.Errorf("failed to create memory: %w", err)
		}
	}

	return memory, nil
}
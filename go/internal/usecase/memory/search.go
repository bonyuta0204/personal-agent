package memory

import (
	"context"
	"fmt"

	"github.com/bonyuta0204/personal-agent/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/internal/domain/port/embedding"
	"github.com/bonyuta0204/personal-agent/internal/domain/port/repository"
)

type SearchUseCase struct {
	memoryRepo      repository.MemoryRepository
	embeddingClient embedding.Provider
}

func NewSearchUseCase(memoryRepo repository.MemoryRepository, embeddingClient embedding.Provider) *SearchUseCase {
	return &SearchUseCase{
		memoryRepo:      memoryRepo,
		embeddingClient: embeddingClient,
	}
}

type SearchByTextRequest struct {
	Query string
	Limit int
}

type SearchByTagsRequest struct {
	Tags []string
}

func (uc *SearchUseCase) SearchByText(ctx context.Context, req *SearchByTextRequest) ([]*model.Memory, error) {
	if req.Limit <= 0 {
		req.Limit = 10 // default limit
	}

	// Generate embedding for the search query
	embeddingVector, err := uc.embeddingClient.GenerateEmbedding(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding for search query: %w", err)
	}

	// Search by embedding similarity
	memories, err := uc.memoryRepo.SearchByEmbedding(ctx, embeddingVector, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories by embedding: %w", err)
	}

	return memories, nil
}

func (uc *SearchUseCase) SearchByTags(ctx context.Context, req *SearchByTagsRequest) ([]*model.Memory, error) {
	if len(req.Tags) == 0 {
		return []*model.Memory{}, nil
	}

	memories, err := uc.memoryRepo.SearchByTags(ctx, req.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories by tags: %w", err)
	}

	return memories, nil
}
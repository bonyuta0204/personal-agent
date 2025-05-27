package openai

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/embedding"
	"github.com/sashabaranov/go-openai"
)

// Provider implements the embedding.EmbeddingProvider interface using OpenAI API
type Provider struct {
	client *openai.Client
	model  openai.EmbeddingModel
}

// NewProvider creates a new OpenAI embedding provider
func NewProvider() (*Provider, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)
	return &Provider{
		client: client,
		model:  openai.AdaEmbeddingV2,
	}, nil
}

// Embed implements the embedding.EmbeddingProvider interface
// It creates an embedding for the given text using OpenAI's API
func (p *Provider) Embed(text string) ([]float64, error) {
	// OpenAI embedding API max: 300,000 tokens (see https://platform.openai.com/docs/guides/embeddings/what-are-embeddings)
	// Roughly estimate: 1 token ≈ 4 chars (so 120,000 chars ≈ 300,000 tokens)
	const maxChars = 120000
	if len([]rune(text)) > maxChars {
		return nil, fmt.Errorf("text too long for embedding: %d chars (max %d chars, ~300k tokens)", len([]rune(text)), maxChars)
	}

	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: p.model,
	}

	resp, err := p.client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no embedding data returned from OpenAI")
	}

	// Convert []float32 to []float64
	embedding := make([]float64, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float64(v)
	}
	return embedding, nil
}

// Ensure Provider implements the EmbeddingProvider interface
var _ embedding.EmbeddingProvider = (*Provider)(nil)

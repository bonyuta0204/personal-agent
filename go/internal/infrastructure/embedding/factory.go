package embedding

import (
	"github.com/bonyuta0204/personal-agent/go/internal/domain/port/embedding"
	"github.com/bonyuta0204/personal-agent/go/internal/infrastructure/embedding/openai"
)

// NewOpenAIProvider creates a new OpenAI embedding provider
func NewOpenAIProvider() (embedding.EmbeddingProvider, error) {
	return openai.NewProvider()
}

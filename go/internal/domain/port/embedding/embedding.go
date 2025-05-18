package embedding

type EmbeddingProvider interface {
	Embed(text string) ([]float64, error)
}

package model

import "time"

type DocumentId string

// represent a document in the knowledge base
type Document struct {
	ID        DocumentId
	StoreId   StoreId
	Path      string
	Content   string
	Embedding []float64
	Tags      []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

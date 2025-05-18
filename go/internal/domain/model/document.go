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
	SHA       string

	ModifiedAt time.Time // The time when the document was last modified. This is used to detect changes in the document.
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// represent a document entry in the knowledge base
type DocumentEntry struct {
	Path       string
	ModifiedAt time.Time
}

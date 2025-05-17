package model

import "time"

type MemoryId string

// represent a memory
type Memory struct {
	ID        MemoryId
	Path      string
	Content   string
	Embedding []float64
	Tags      []string

	ModifiedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

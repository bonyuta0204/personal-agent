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

	CreatedAt time.Time
	UpdatedAt time.Time
}

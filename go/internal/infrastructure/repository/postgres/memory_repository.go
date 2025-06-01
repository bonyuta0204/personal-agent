package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	repo "github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Ensure memoryRepository implements repo.MemoryRepository
var _ repo.MemoryRepository = (*memoryRepository)(nil)

type memoryRepository struct {
	db *sqlx.DB
}

// NewMemoryRepository creates a new PostgreSQL memory repository
func NewMemoryRepository(db *sqlx.DB) repo.MemoryRepository {
	return &memoryRepository{db: db}
}

// SaveMemory saves or updates a memory in the database
func (r *memoryRepository) SaveMemory(memory *model.Memory) error {
	if memory == nil {
		return errors.New("memory cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Convert tags to JSONB (always as array, never null)
	tags := memory.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	// Convert embedding to PostgreSQL vector format
	const expectedDim = 1536
	var embeddingStr string = ""
	if len(memory.Embedding) > 0 {
		if len(memory.Embedding) != expectedDim {
			return fmt.Errorf("invalid embedding dimension: got %d, want %d", len(memory.Embedding), expectedDim)
		}
		for i, v := range memory.Embedding {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				return fmt.Errorf("invalid embedding value at position %d: %v", i, v)
			}
		}
		// Convert to JSON array format
		parts := make([]string, len(memory.Embedding))
		for i, v := range memory.Embedding {
			parts[i] = fmt.Sprintf("%f", v)
		}
		embeddingStr = "[" + strings.Join(parts, ",") + "]"
	}

	// Check if memory exists
	var exists bool
	err = tx.GetContext(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM memories WHERE path = $1)`,
		memory.Path,
	)
	if err != nil {
		return err
	}

	if exists {
		// Update existing memory
		_, err = tx.ExecContext(ctx, `
			UPDATE memories 
			SET content = $1, 
			    embedding = $2, 
			    tags = $3,
			    modified_at = $4,
			    sha = $5,
			    updated_at = NOW()
			WHERE path = $6`,
			memory.Content,
			embeddingStr,
			tagsJSON,
			memory.ModifiedAt,
			memory.SHA,
			memory.Path,
		)
	} else {
		// Insert new memory
		_, err = tx.ExecContext(ctx, `
			INSERT INTO memories (path, content, embedding, tags, modified_at, sha)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			memory.Path,
			memory.Content,
			embeddingStr,
			tagsJSON,
			memory.ModifiedAt,
			memory.SHA,
		)
	}

	if err != nil {
		return err
	}

	// Get the updated/inserted memory to set timestamps
	var updatedMem struct {
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	err = tx.GetContext(ctx, &updatedMem,
		`SELECT created_at, updated_at FROM memories WHERE path = $1`,
		memory.Path,
	)
	if err != nil {
		return err
	}

	memory.CreatedAt = updatedMem.CreatedAt
	memory.UpdatedAt = updatedMem.UpdatedAt

	return tx.Commit()
}

// ListMemories retrieves all memories from the database
func (r *memoryRepository) ListMemories() ([]*model.Memory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT 
			id,
			path,
			content,
			embedding,
			tags,
			sha,
			modified_at,
			created_at,
			updated_at
		FROM memories
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query memories: %w", err)
	}
	defer rows.Close()

	var memories []*model.Memory
	for rows.Next() {
		var memory model.Memory
		var embeddingStr string
		var tagsJSON []byte

		err := rows.Scan(
			&memory.ID,
			&memory.Path,
			&memory.Content,
			&embeddingStr,
			&tagsJSON,
			&memory.SHA,
			&memory.ModifiedAt,
			&memory.CreatedAt,
			&memory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory row: %w", err)
		}

		// Parse tags from JSON
		if len(tagsJSON) > 0 {
			if err := json.Unmarshal(tagsJSON, &memory.Tags); err != nil {
				return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
			}
		}

		// Parse embedding from vector string format
		if embeddingStr != "" && embeddingStr != "[]" {
			// Remove brackets and split by comma
			embeddingStr = strings.Trim(embeddingStr, "[]")
			parts := strings.Split(embeddingStr, ",")
			memory.Embedding = make([]float64, len(parts))
			for i, part := range parts {
				var val float64
				_, err := fmt.Sscanf(strings.TrimSpace(part), "%f", &val)
				if err != nil {
					return nil, fmt.Errorf("failed to parse embedding value: %w", err)
				}
				memory.Embedding[i] = val
			}
		}

		memories = append(memories, &memory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating memory rows: %w", err)
	}

	return memories, nil
}

// FindExistingSHAs returns the SHAs of memories that exist in the database
// This is used to find unchanged memories that don't need to be updated
func (r *memoryRepository) FindExistingSHAs(memories []*model.Memory) ([]string, error) {
	if len(memories) == 0 {
		return nil, nil
	}

	var shas []string

	for _, mem := range memories {
		if mem != nil {
			shas = append(shas, mem.SHA)
		}
	}

	if len(shas) == 0 {
		return nil, nil
	}

	// Get existing memories with their SHAs
	query := `
		SELECT sha
		FROM memories
		WHERE sha = ANY($1)
	`

	rows, err := r.db.Query(query, pq.Array(shas))
	if err != nil {
		return nil, fmt.Errorf("failed to query memories: %w", err)
	}
	defer rows.Close()

	// Collect memory SHAs where SHA matches
	var unchangedSHAs []string

	for rows.Next() {
		var sha string

		if err := rows.Scan(&sha); err != nil {
			return nil, fmt.Errorf("failed to scan memory row: %w", err)
		}

		unchangedSHAs = append(unchangedSHAs, sha)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating memory rows: %w", err)
	}

	return unchangedSHAs, nil
}
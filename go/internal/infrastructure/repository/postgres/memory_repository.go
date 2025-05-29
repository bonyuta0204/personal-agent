package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/bonyuta0204/personal-agent/internal/domain/model"
	"github.com/bonyuta0204/personal-agent/internal/domain/port/repository"
	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type memoryRepository struct {
	db *sql.DB
}

func NewMemoryRepository(db *sql.DB) repository.MemoryRepository {
	return &memoryRepository{db: db}
}

func (r *memoryRepository) Create(ctx context.Context, memory *model.Memory) (*model.Memory, error) {
	query := `
		INSERT INTO memories (path, content, embedding, tags, sha, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	tagsJSON, err := json.Marshal(memory.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}

	var embedding pgvector.Vector
	if len(memory.Embedding) > 0 {
		embedding = pgvector.NewVector(memory.Embedding)
	}

	var id int
	err = r.db.QueryRowContext(
		ctx,
		query,
		memory.Path,
		memory.Content,
		embedding,
		tagsJSON,
		memory.SHA,
		memory.ModifiedAt,
	).Scan(&id, &memory.CreatedAt, &memory.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create memory: %w", err)
	}

	memory.ID = model.MemoryId(fmt.Sprintf("%d", id))
	return memory, nil
}

func (r *memoryRepository) GetByID(ctx context.Context, id model.MemoryId) (*model.Memory, error) {
	query := `
		SELECT id, path, content, embedding, tags, sha, modified_at, created_at, updated_at
		FROM memories
		WHERE id = $1
	`

	var memory model.Memory
	var idInt int
	var embedding pgvector.Vector
	var tagsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&idInt,
		&memory.Path,
		&memory.Content,
		&embedding,
		&tagsJSON,
		&memory.SHA,
		&memory.ModifiedAt,
		&memory.CreatedAt,
		&memory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("memory not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	memory.ID = model.MemoryId(fmt.Sprintf("%d", idInt))
	memory.Embedding = embedding.Slice()

	if err := json.Unmarshal(tagsJSON, &memory.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &memory, nil
}

func (r *memoryRepository) GetByPath(ctx context.Context, path string) (*model.Memory, error) {
	query := `
		SELECT id, path, content, embedding, tags, sha, modified_at, created_at, updated_at
		FROM memories
		WHERE path = $1
	`

	var memory model.Memory
	var idInt int
	var embedding pgvector.Vector
	var tagsJSON []byte

	err := r.db.QueryRowContext(ctx, query, path).Scan(
		&idInt,
		&memory.Path,
		&memory.Content,
		&embedding,
		&tagsJSON,
		&memory.SHA,
		&memory.ModifiedAt,
		&memory.CreatedAt,
		&memory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get memory by path: %w", err)
	}

	memory.ID = model.MemoryId(fmt.Sprintf("%d", idInt))
	memory.Embedding = embedding.Slice()

	if err := json.Unmarshal(tagsJSON, &memory.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return &memory, nil
}

func (r *memoryRepository) List(ctx context.Context, offset, limit int) ([]*model.Memory, error) {
	query := `
		SELECT id, path, content, embedding, tags, sha, modified_at, created_at, updated_at
		FROM memories
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list memories: %w", err)
	}
	defer rows.Close()

	var memories []*model.Memory
	for rows.Next() {
		var memory model.Memory
		var idInt int
		var embedding pgvector.Vector
		var tagsJSON []byte

		err := rows.Scan(
			&idInt,
			&memory.Path,
			&memory.Content,
			&embedding,
			&tagsJSON,
			&memory.SHA,
			&memory.ModifiedAt,
			&memory.CreatedAt,
			&memory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory: %w", err)
		}

		memory.ID = model.MemoryId(fmt.Sprintf("%d", idInt))
		memory.Embedding = embedding.Slice()

		if err := json.Unmarshal(tagsJSON, &memory.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		memories = append(memories, &memory)
	}

	return memories, nil
}

func (r *memoryRepository) Update(ctx context.Context, memory *model.Memory) error {
	query := `
		UPDATE memories
		SET path = $2, content = $3, embedding = $4, tags = $5, sha = $6, modified_at = $7
		WHERE id = $1
	`

	tagsJSON, err := json.Marshal(memory.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	var embedding pgvector.Vector
	if len(memory.Embedding) > 0 {
		embedding = pgvector.NewVector(memory.Embedding)
	}

	result, err := r.db.ExecContext(
		ctx,
		query,
		memory.ID,
		memory.Path,
		memory.Content,
		embedding,
		tagsJSON,
		memory.SHA,
		memory.ModifiedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update memory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("memory not found: %s", memory.ID)
	}

	return nil
}

func (r *memoryRepository) Delete(ctx context.Context, id model.MemoryId) error {
	query := `DELETE FROM memories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("memory not found: %s", id)
	}

	return nil
}

func (r *memoryRepository) SearchByEmbedding(ctx context.Context, embedding []float64, limit int) ([]*model.Memory, error) {
	query := `
		SELECT id, path, content, embedding, tags, sha, modified_at, created_at, updated_at,
		       1 - (embedding <=> $1::vector) as similarity
		FROM memories
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> $1::vector
		LIMIT $2
	`

	embeddingVector := pgvector.NewVector(embedding)
	rows, err := r.db.QueryContext(ctx, query, embeddingVector, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories by embedding: %w", err)
	}
	defer rows.Close()

	var memories []*model.Memory
	for rows.Next() {
		var memory model.Memory
		var idInt int
		var embeddingResult pgvector.Vector
		var tagsJSON []byte
		var similarity float64

		err := rows.Scan(
			&idInt,
			&memory.Path,
			&memory.Content,
			&embeddingResult,
			&tagsJSON,
			&memory.SHA,
			&memory.ModifiedAt,
			&memory.CreatedAt,
			&memory.UpdatedAt,
			&similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory: %w", err)
		}

		memory.ID = model.MemoryId(fmt.Sprintf("%d", idInt))
		memory.Embedding = embeddingResult.Slice()

		if err := json.Unmarshal(tagsJSON, &memory.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		memories = append(memories, &memory)
	}

	return memories, nil
}

func (r *memoryRepository) SearchByTags(ctx context.Context, tags []string) ([]*model.Memory, error) {
	query := `
		SELECT id, path, content, embedding, tags, sha, modified_at, created_at, updated_at
		FROM memories
		WHERE tags @> $1
		ORDER BY created_at DESC
	`

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search tags: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories by tags: %w", err)
	}
	defer rows.Close()

	var memories []*model.Memory
	for rows.Next() {
		var memory model.Memory
		var idInt int
		var embedding pgvector.Vector
		var tagsJSON []byte

		err := rows.Scan(
			&idInt,
			&memory.Path,
			&memory.Content,
			&embedding,
			&tagsJSON,
			&memory.SHA,
			&memory.ModifiedAt,
			&memory.CreatedAt,
			&memory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory: %w", err)
		}

		memory.ID = model.MemoryId(fmt.Sprintf("%d", idInt))
		memory.Embedding = embedding.Slice()

		if err := json.Unmarshal(tagsJSON, &memory.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		memories = append(memories, &memory)
	}

	return memories, nil
}
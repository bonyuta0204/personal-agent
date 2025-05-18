package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/bonyuta0204/personal-agent/go/internal/domain/model"
	repo "github.com/bonyuta0204/personal-agent/go/internal/domain/port/repository"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Ensure documentRepository implements repo.DocumentRepository
var _ repo.DocumentRepository = (*documentRepository)(nil)

type documentRepository struct {
	db *sqlx.DB
}

// NewDocumentRepository creates a new PostgreSQL document repository
func NewDocumentRepository(db *sqlx.DB) repo.DocumentRepository {
	return &documentRepository{db: db}
}

// SaveDocument saves or updates a document in the database
func (r *documentRepository) SaveDocument(document *model.Document) error {
	if document == nil {
		return errors.New("document cannot be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Convert tags to JSONB
	tagsJSON, err := json.Marshal(document.Tags)
	if err != nil {
		return err
	}

	// Convert embedding to PostgreSQL vector format
	var embedding interface{} = nil
	if len(document.Embedding) > 0 {
		sql := "SELECT $1::vector"
		err = tx.Get(&embedding, sql, pq.Array(document.Embedding))
		if err != nil {
			return err
		}
	}

	// Check if document exists
	var exists bool
	err = tx.GetContext(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM documents WHERE store_id = $1 AND path = $2)`,
		document.StoreId, document.Path,
	)
	if err != nil {
		return err
	}

	log.Printf("Saving document %s", document.Path)
	if exists {
		// Update existing document
		_, err = tx.ExecContext(ctx, `
			UPDATE documents 
			SET content = $1, 
			    embedding = $2, 
			    tags = $3,
			    modified_at = $4,
			    sha = $5,
			    updated_at = NOW()
			WHERE store_id = $6 AND path = $7`,
			document.Content,
			embedding,
			tagsJSON,
			document.ModifiedAt,
			document.SHA,
			document.StoreId,
			document.Path,
		)
	} else {
		// Insert new document
		_, err = tx.ExecContext(ctx, `
			INSERT INTO documents (store_id, path, content, embedding, tags, modified_at, sha)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
			document.StoreId,
			document.Path,
			document.Content,
			embedding,
			tagsJSON,
			document.ModifiedAt,
			document.SHA,
		)
	}

	if err != nil {
		return err
	}

	// Get the updated/inserted document to set timestamps
	var updatedDoc struct {
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	err = tx.GetContext(ctx, &updatedDoc,
		`SELECT created_at, updated_at FROM documents WHERE store_id = $1 AND path = $2`,
		document.StoreId, document.Path,
	)
	if err != nil {
		return err
	}

	document.CreatedAt = updatedDoc.CreatedAt
	document.UpdatedAt = updatedDoc.UpdatedAt

	return tx.Commit()
}

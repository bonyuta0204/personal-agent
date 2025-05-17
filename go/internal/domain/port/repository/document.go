package repository

import "github.com/bonyuta0204/personal-agent/go/internal/domain/model"

type DocumentRepository interface {
	SaveDocument(document *model.Document) error
}

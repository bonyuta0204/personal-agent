package repository

import "github.com/bonyuta0204/personal-agent/go/internal/domain/model"

type DocumentRepository interface {
	SaveDocument(document *model.Document) error
	// FindExistingSHAs returns the IDs of documents that are unchanged
	FindExistingSHAs(documents []*model.Document) ([]string, error)
}

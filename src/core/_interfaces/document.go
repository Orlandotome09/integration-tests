package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type DocumentAdapter interface {
	GetByID(id string) (*entity.Document, error)
	FindByEntityIDAndDocumentType(id string, documentType string) ([]entity.Document, error)
	Find(entityID string) ([]entity.Document, error)
}

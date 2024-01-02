package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type DocumentFileAdapter interface {
	Get(documentFileID uuid.UUID) (*entity.DocumentFile, error)
	FindByDocumentID(documentID uuid.UUID) ([]entity.DocumentFile, error)
	GetLastTwoFilesOfDocument(documentID uuid.UUID) (*entity.DocumentFile, *entity.DocumentFile, error)
}

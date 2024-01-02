package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type DOAResultRepository interface {
	Get(id *uuid.UUID) (*entity.DOAResult, error)
	Save(doaResult *entity.DOAResult) (*entity.DOAResult, error)
	Enrich(doaResult *entity.DOAResult) (*entity.DOAResult, error)
	FindByEntityID(entityID *uuid.UUID) ([]entity.DOAResult, error)
	FindByEntityIdAndDocumentId(entityID *uuid.UUID, documentID *uuid.UUID) ([]entity.DOAResult, error)
	FindLastByEntityIdAndDocumentId(entityID *uuid.UUID, documentID *uuid.UUID) (*entity.DOAResult, error)
}

package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa/contracts"
	"github.com/google/uuid"
)

type DOAService interface {
	Get(id *uuid.UUID) (*entity2.DOAResult, error)
	Save(doaResult *entity2.DOAResult) (*entity2.DOAResult, error)
	EnrichWithScores(id *uuid.UUID, scores entity2.Scores) (*entity2.DOAResult, error)
	FindLastResult(entityID *uuid.UUID, documentID *uuid.UUID) (*entity2.DOAResult, error)
	RequestExtraction(frontFile *entity2.DocumentFile, frontFileURI string,
		backFile *entity2.DocumentFile, backFileURI string,
		doc *entity2.Document, profileID uuid.UUID) (*entity2.DOAExtraction, error)
	CreateDocMetadata(doaResult *entity2.DOAResult, documentType string,
		docSubType string, documentID string, profileID string) *contracts.DocumentMetadata
}

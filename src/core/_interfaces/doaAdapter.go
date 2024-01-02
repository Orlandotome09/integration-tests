package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type DOAAdapter interface {
	RequestExtraction(frontFile *entity2.DocumentFile, frontFileURI string,
		backFile *entity2.DocumentFile, backFileURI string, doc *entity2.Document,
		profileID uuid.UUID) (*entity2.DOAExtraction, error)
}

package doa

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa/contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockDOAService struct {
	_interfaces.DOAService
	mock.Mock
}

func (ref *MockDOAService) Get(id *uuid.UUID) (*entity2.DOAResult, error) {
	args := ref.Called(id)
	return args.Get(0).(*entity2.DOAResult), args.Error(1)
}

func (ref *MockDOAService) Save(doaResult *entity2.DOAResult) (*entity2.DOAResult, error) {
	args := ref.Called(doaResult)
	return args.Get(0).(*entity2.DOAResult), args.Error(1)
}

func (ref *MockDOAService) EnrichWithScores(id *uuid.UUID, scores entity2.Scores) (*entity2.DOAResult, error) {
	args := ref.Called(id, scores)
	return args.Get(0).(*entity2.DOAResult), args.Error(1)
}

func (ref *MockDOAService) FindLastResult(entityID *uuid.UUID, documentID *uuid.UUID) (*entity2.DOAResult, error) {
	args := ref.Called(entityID, documentID)
	return args.Get(0).(*entity2.DOAResult), args.Error(1)
}

func (ref *MockDOAService) RequestExtraction(frontFile *entity2.DocumentFile, frontFileURI string,
	backFile *entity2.DocumentFile, backFileURI string,
	doc *entity2.Document, profileID uuid.UUID) (*entity2.DOAExtraction, error) {
	args := ref.Called(frontFile, frontFileURI, backFile, backFileURI, doc, profileID)
	return args.Get(0).(*entity2.DOAExtraction), args.Error(1)
}

func (ref *MockDOAService) CreateDocMetadata(doaResult *entity2.DOAResult, documentType string,
	docSubType string, documentID string, profileID string) *contracts.DocumentMetadata {
	args := ref.Called(doaResult, documentType, docSubType, documentID, profileID)
	return args.Get(0).(*contracts.DocumentMetadata)
}

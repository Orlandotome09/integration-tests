package doaTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/doa/contracts"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockDOATranslator struct {
	DOATranslator
	mock.Mock
}

func (ref *MockDOATranslator) FromDomain(frontFile *entity2.DocumentFile, frontFileURI string,
	backFile *entity2.DocumentFile, backFileURI string, doc *entity2.Document,
	profileID uuid.UUID) (*contracts.DOAExtractionRequest, error) {
	args := ref.Called(frontFile, frontFileURI, backFile, backFileURI, doc, profileID)
	return args.Get(0).(*contracts.DOAExtractionRequest), args.Error(1)
}

func (ref *MockDOATranslator) ToDomain(response []byte) (*entity2.DOAExtraction, error) {
	args := ref.Called(response)
	return args.Get(0).(*entity2.DOAExtraction), args.Error(1)
}

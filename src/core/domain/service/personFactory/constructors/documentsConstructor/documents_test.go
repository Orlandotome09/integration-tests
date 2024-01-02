package documentsConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_When_Incomplete_Documents_Is_Enabled_Should_Get_Documents(t *testing.T) {
	documentAdapter := &mocks.DocumentAdapter{}
	documentFileAdapter := &mocks.DocumentFileAdapter{}
	constructor := documentsPersonConstructor{documentAdapter: documentAdapter,
		documentFileAdapter: documentFileAdapter}

	person := entity.Person{
		EntityID: uuid.New(),
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						IncompleteParams: &entity.IncompleteParams{DocumentsRequired: []entity.DocumentRequired{{DocumentType: values.DocumentTypeIdentification}}},
					},
				},
			},
		},
	}
	personWrapper := entity.PersonWrapper{
		Person: person,
	}
	documents := []entity.Document{{DocumentID: uuid.New()}}
	documentFileID := uuid.New()
	documentFiles := []entity.DocumentFile{{DocumentFileID: &documentFileID}}

	documentAdapter.On("Find", person.EntityID.String()).Return(documents, nil)
	documentFileAdapter.On("FindByDocumentID", documents[0].DocumentID).Return(documentFiles, nil)

	err := constructor.Assemble(&personWrapper)

	expectedDocuments := documents
	expectedDocumentFiles := documentFiles

	assert.Nil(t, err)
	assert.Equal(t, expectedDocuments, personWrapper.Person.Documents)
	assert.Equal(t, expectedDocumentFiles, personWrapper.Person.DocumentFiles)
	mock.AssertExpectationsForObjects(t, documentAdapter)
}

func Test_Assemble_When_DOA_Rule_Is_Enabled_Should_Get_Documents(t *testing.T) {
	documentAdapter := &mocks.DocumentAdapter{}
	documentFileAdapter := &mocks.DocumentFileAdapter{}
	constructor := documentsPersonConstructor{documentAdapter: documentAdapter,
		documentFileAdapter: documentFileAdapter}

	person := entity.Person{
		EntityID: uuid.New(),
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						DOAParams: &entity.DOAParams{},
					},
				},
			},
		},
	}
	personWrapper := entity.PersonWrapper{
		Person: person,
	}
	documents := []entity.Document{{DocumentID: uuid.New()}}
	documentFileID := uuid.New()
	documentFiles := []entity.DocumentFile{{DocumentFileID: &documentFileID}}

	documentAdapter.On("Find", person.EntityID.String()).Return(documents, nil)
	documentFileAdapter.On("FindByDocumentID", documents[0].DocumentID).Return(documentFiles, nil)

	err := constructor.Assemble(&personWrapper)

	expectedDocuments := documents
	expectedDocumentFiles := documentFiles

	assert.Nil(t, err)
	assert.Equal(t, expectedDocuments, personWrapper.Person.Documents)
	assert.Equal(t, expectedDocumentFiles, personWrapper.Person.DocumentFiles)
	mock.AssertExpectationsForObjects(t, documentAdapter)
}

func Test_Assemble_Should_Not_Get_Documents(t *testing.T) {
	documentAdapter := &mocks.DocumentAdapter{}
	constructor := documentsPersonConstructor{documentAdapter: documentAdapter}

	person := entity.Person{
		EntityID: uuid.New(),
	}
	personWrapper := entity.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(&personWrapper)

	var expected []entity.Document = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.Documents)
	mock.AssertExpectationsForObjects(t, documentAdapter)
}

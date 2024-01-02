package documentTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type DocumentTranslator interface {
	Translate(response contracts.DocumentResponse) entity.Document
	TranslateAll(responses []contracts.DocumentResponse) []entity.Document
}

type documentFileTranslator struct{}

func New() DocumentTranslator {
	return &documentFileTranslator{}
}

func (ref *documentFileTranslator) Translate(response contracts.DocumentResponse) entity.Document {
	document := entity.Document{
		DocumentID:      response.DocumentID,
		EntityID:        response.EntityID,
		DocumentType:    response.Type,
		DocumentSubType: response.SubType,
		DocumentFields: entity.DocumentFields{
			Number:    response.DocumentFields.Number,
			IssueDate: response.DocumentFields.IssueDate,
			Name:      response.DocumentFields.Name,
		},
		ExpirationDate: response.ExpirationDate,
		EmissionDate:   response.EmissionDate,
	}
	return document
}

func (ref *documentFileTranslator) TranslateAll(responses []contracts.DocumentResponse) []entity.Document {
	documents := []entity.Document{}

	for _, elem := range responses {
		document := ref.Translate(elem)
		documents = append(documents, document)
	}

	return documents
}

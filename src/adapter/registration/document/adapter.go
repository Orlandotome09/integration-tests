package document

import (
	documentClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http"
	documentTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type documentAdapter struct {
	documentClient documentClient.DocumentClient
	translator     documentTranslator.DocumentTranslator
}

func NewDocumentAdapter(documentClient documentClient.DocumentClient,
	translator documentTranslator.DocumentTranslator) interfaces.DocumentAdapter {
	return &documentAdapter{
		documentClient: documentClient,
		translator:     translator,
	}
}

func (ref *documentAdapter) GetByID(id string) (*entity.Document, error) {

	resp, err := ref.documentClient.Get(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, nil
	}

	document := ref.translator.Translate(*resp)

	return &document, nil
}

func (ref *documentAdapter) FindByEntityIDAndDocumentType(id string,
	documentType string) ([]entity.Document, error) {

	responses, err := ref.documentClient.SearchByEntityID(id)
	if err != nil {
		return []entity.Document{}, errors.WithStack(err)
	}

	documents := ref.translator.TranslateAll(responses)

	documents = filterByDocumentType(documentType, documents)

	return documents, nil
}

func filterByDocumentType(documentType string, documents []entity.Document) []entity.Document {

	filtered := []entity.Document{}

	for _, document := range documents {
		if document.DocumentType == documentType {
			filtered = append(filtered, document)
		}
	}

	return filtered
}

func (ref *documentAdapter) Find(entityID string) ([]entity.Document, error) {

	responses, err := ref.documentClient.SearchByEntityID(entityID)
	if err != nil {
		return []entity.Document{}, errors.WithStack(err)
	}

	documents := ref.translator.TranslateAll(responses)

	return documents, nil
}

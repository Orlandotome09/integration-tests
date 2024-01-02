package doaAdapter

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter"
	translator "bitbucket.org/bexstech/temis-compliance/src/adapter/doa/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

const (
	GetRequestExtractionPath = "/extraction/"
)

type doaAdapter struct {
	doaHttpClient adapter.HttpClient
	doaTranslator translator.DOATranslator
}

func NewDOAAdapter(doaHttpClient adapter.HttpClient, doaTranslator translator.DOATranslator) interfaces.DOAAdapter {
	return &doaAdapter{
		doaHttpClient: doaHttpClient,
		doaTranslator: doaTranslator,
	}
}

func (ref *doaAdapter) RequestExtraction(frontFile *entity.DocumentFile, frontFileURI string,
	backFile *entity.DocumentFile, backFileURI string, doc *entity.Document,
	profileID uuid.UUID) (*entity.DOAExtraction, error) {

	doaExtractionRequest, err := ref.doaTranslator.FromDomain(frontFile, frontFileURI, backFile, backFileURI,
		doc, profileID)
	if err != nil {
		return nil, err
	}

	response, err := ref.doaHttpClient.Post(GetRequestExtractionPath, doaExtractionRequest)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}

	return ref.doaTranslator.ToDomain(response)
}

package doaTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/doa/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DOATranslator interface {
	FromDomain(frontFile *entity.DocumentFile, frontFileURI string,
		backFile *entity.DocumentFile, backFileURI string, doc *entity.Document,
		profileID uuid.UUID) (*contracts.DOAExtractionRequest, error)
	ToDomain(response []byte) (*entity.DOAExtraction, error)
}

type doaTranslator struct {
	temisComplianceHost string
}

func New(temisComplianceHost string) DOATranslator {
	return &doaTranslator{
		temisComplianceHost: temisComplianceHost,
	}
}

func (ref *doaTranslator) FromDomain(frontFile *entity.DocumentFile, frontFileURI string,
	backFile *entity.DocumentFile, backFileURI string, doc *entity.Document,
	profileID uuid.UUID) (*contracts.DOAExtractionRequest, error) {

	file1 := contracts.FileParams{
		FileID:   frontFile.FileID.String(),
		FileSide: contracts.ToFileSide(frontFile.FileSide),
		FileURI:  frontFileURI,
		FileName: frontFile.FileID.String(),
	}

	file2 := contracts.FileParams{
		FileID:   backFile.FileID.String(),
		FileSide: contracts.ToFileSide(backFile.FileSide),
		FileURI:  backFileURI,
		FileName: backFile.FileID.String(),
	}

	requestBody := &contracts.DOAExtractionRequest{
		ProfileID:      profileID,
		Type:           contracts.ToDocumentType(doc.DocumentSubType),
		Metadata:       contracts.ToMetadata(doc.DocumentFields),
		FileParams:     []contracts.FileParams{file1, file2},
		CallbackParams: &contracts.CallbackParams{URL: ref.temisComplianceHost + "/doa/callback/"},
	}

	return requestBody, nil
}

func (ref *doaTranslator) ToDomain(response []byte) (*entity.DOAExtraction, error) {
	var doaResponse contracts.DOAExtractionResponse
	err := json.Unmarshal(response, &doaResponse)
	if err != nil {
		return nil, errors.Errorf("could not convert, response is not a doa response: %+v", string(response))
	}

	return &entity.DOAExtraction{
		Message:   doaResponse.Message,
		RequestID: doaResponse.RequestID,
	}, nil
}

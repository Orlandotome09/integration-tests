package documentFileTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type DocumentFileTranslator interface {
	Translate(response contracts.DocumentFileResponse) entity.DocumentFile
	TranslateAll(responses []contracts.DocumentFileResponse) []entity.DocumentFile
}

type documentFileTranslator struct{}

func New() DocumentFileTranslator {
	return &documentFileTranslator{}
}

func (ref *documentFileTranslator) Translate(response contracts.DocumentFileResponse) entity.DocumentFile {
	documentFile := entity.DocumentFile{
		DocumentFileID: &response.DocumentFileID,
		DocumentID:     response.DocumentID,
		FileID:         response.FileID,
		FileSide:       values.FileSide(response.FileSide),
		CreatedAt:      response.CreatedAt,
	}

	return documentFile
}

func (ref *documentFileTranslator) TranslateAll(responses []contracts.DocumentFileResponse) []entity.DocumentFile {
	documentFiles := []entity.DocumentFile{}

	for _, elem := range responses {
		documentFile := ref.Translate(elem)
		documentFiles = append(documentFiles, documentFile)
	}

	return documentFiles
}

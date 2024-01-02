package documentFileTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockDocumentFileTranslator struct {
	DocumentFileTranslator
	mock.Mock
}

func (ref *MockDocumentFileTranslator) Translate(response contracts.DocumentFileResponse) entity.DocumentFile {
	args := ref.Called(response)
	return args.Get(0).(entity.DocumentFile)
}

func (ref *MockDocumentFileTranslator) TranslateAll(responses []contracts.DocumentFileResponse) []entity.DocumentFile {
	args := ref.Called(responses)
	return args.Get(0).([]entity.DocumentFile)
}

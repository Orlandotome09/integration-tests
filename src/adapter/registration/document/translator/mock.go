package documentTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockDocumentTranslator struct {
	DocumentTranslator
	mock.Mock
}

func (ref *MockDocumentTranslator) Translate(response contracts.DocumentResponse) entity.Document {
	args := ref.Called(response)
	return args.Get(0).(entity.Document)
}

func (ref *MockDocumentTranslator) TranslateAll(responses []contracts.DocumentResponse) []entity.Document {
	args := ref.Called(responses)
	return args.Get(0).([]entity.Document)
}

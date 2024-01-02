package translator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockContactTranslator struct {
	Translator
	mock.Mock
}

func (ref *MockContactTranslator) Translate(responses []contracts.ContactResponse) []entity.Contact {
	args := ref.Called(responses)
	return args.Get(0).([]entity.Contact)
}

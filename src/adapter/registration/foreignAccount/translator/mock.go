package foreignAccountTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockForeignAccountTranslator struct {
	ForeignAccountTranslator
	mock.Mock
}

func (ref *MockForeignAccountTranslator) Translate(response contracts.ForeignAccountResponse) entity.ForeignAccount {
	args := ref.Called(response)
	return args.Get(0).(entity.ForeignAccount)
}

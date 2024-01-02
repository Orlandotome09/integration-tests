package registrationAddressesTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockAddressTranslator struct {
	Translator
	mock.Mock
}

func (ref *MockAddressTranslator) Translate(responses []contracts.AddressResponse) []entity.Address {
	args := ref.Called(responses)
	return args.Get(0).([]entity.Address)
}

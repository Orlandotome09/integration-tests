package translator

import (
	manuallyOwnershipStructureHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockOwnershipStructureTranslator struct {
	OwnershipStructureTranslator
	mock.Mock
}

func (ref *MockOwnershipStructureTranslator) Translate(response manuallyOwnershipStructureHttpClient.OwnershipStructureResponse) (*entity.OwnershipStructure, error) {
	args := ref.Called(response)
	return args.Get(0).(*entity.OwnershipStructure), args.Error(1)
}

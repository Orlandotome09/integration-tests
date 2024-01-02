package legalRepresentativeTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockLegalRepresentativeTranslator struct {
	LegalRepresentativeTranslator
	mock.Mock
}

func (ref *MockLegalRepresentativeTranslator) Translate(response contracts.LegalRepresentativeResponse) (*entity.LegalRepresentative, error) {
	args := ref.Called(response)
	return args.Get(0).(*entity.LegalRepresentative), args.Error(1)
}

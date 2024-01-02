package enrichedLegalEntityTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/_interfacesEnrichment"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockEnrichedLegalEntityTranslator struct {
	_interfacesEnrichment.EnrichedInformationTranslator
	mock.Mock
}

func (ref *MockEnrichedLegalEntityTranslator) Translate(response []byte) (*entity.EnrichedInformation, error) {
	args := ref.Called(response)
	return args.Get(0).(*entity.EnrichedInformation), args.Error(1)
}

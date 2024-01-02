package enricherTranslator

import (
	contracts "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockEnricherTranslator struct {
	EnricherTranslator
	mock.Mock
}

func (ref *MockEnricherTranslator) Translate(response contracts.EnricherResponse) (*entity.EnrichedInformation, error) {
	args := ref.Called(response)
	return args.Get(0).(*entity.EnrichedInformation), args.Error(1)
}

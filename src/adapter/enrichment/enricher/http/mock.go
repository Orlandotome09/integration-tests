package enricherClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockEnricherClient struct {
	EnricherClient
	mock.Mock
}

func (ref *MockEnricherClient) GetEnrichedPerson(request enricherContracts.EnricherRequest, documentNumber string) (*enricherContracts.EnricherResponse, error) {
	args := ref.Called(request, documentNumber)
	return args.Get(0).(*enricherContracts.EnricherResponse), args.Error(1)
}

package partnerClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockPartnerClient struct {
	PartnerClient
	mock.Mock
}

func (ref *MockPartnerClient) Get(id string) (*contracts.PartnerResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.PartnerResponse), args.Error(1)
}

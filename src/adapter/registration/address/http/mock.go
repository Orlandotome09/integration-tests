package addressClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockAddressClient struct {
	AddressClient
	mock.Mock
}

func (ref *MockAddressClient) Get(id string) (*contracts.AddressResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.AddressResponse), args.Error(1)
}

func (ref *MockAddressClient) Search(profileID string) ([]contracts.AddressResponse, error) {
	args := ref.Called(profileID)
	return args.Get(0).([]contracts.AddressResponse), args.Error(1)
}

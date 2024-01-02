package accountClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockAccountClient struct {
	AccountClient
	mock.Mock
}

func (ref *MockAccountClient) Get(id string) (*contracts.AccountResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.AccountResponse), args.Error(1)
}

func (ref *MockAccountClient) Find(profileID string) ([]contracts.AccountResponse, error) {
	args := ref.Called(profileID)
	return args.Get(0).([]contracts.AccountResponse), args.Error(1)
}

func (ref *MockAccountClient) CreateInternal(request *contracts.CreateAccountRequest) (*contracts.AccountResponse, error) {
	args := ref.Called(request)
	return args.Get(0).(*contracts.AccountResponse), args.Error(1)
}

package contractClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http/contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockContractClient struct {
	ContractClient
	mock.Mock
}

func (ref *MockContractClient) Get(id *uuid.UUID) (*contracts.GetContractResponse, bool, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.GetContractResponse), args.Get(1).(bool), args.Error(2)
}

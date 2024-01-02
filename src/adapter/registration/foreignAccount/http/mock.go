package foreignAccountClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockForeignAccountClient struct {
	ForeignAccountClient
	mock.Mock
}

func (ref *MockForeignAccountClient) Get(id string) (*contracts.ForeignAccountResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.ForeignAccountResponse), args.Error(1)
}

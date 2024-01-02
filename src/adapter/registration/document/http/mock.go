package documentClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockDocumentClient struct {
	DocumentClient
	mock.Mock
}

func (ref *MockDocumentClient) Get(id string) (*contracts.DocumentResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.DocumentResponse), args.Error(1)
}

func (ref *MockDocumentClient) SearchByEntityID(id string) ([]contracts.DocumentResponse, error) {
	args := ref.Called(id)
	return args.Get(0).([]contracts.DocumentResponse), args.Error(1)
}

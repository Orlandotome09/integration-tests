package documentFileClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockDocumentFileClient struct {
	DocumentFileClient
	mock.Mock
}

func (ref *MockDocumentFileClient) Get(documentFileID string) (*contracts.DocumentFileResponse, error) {
	args := ref.Called(documentFileID)
	return args.Get(0).(*contracts.DocumentFileResponse), args.Error(1)
}

func (ref *MockDocumentFileClient) SearchByDocumentId(id string) ([]contracts.DocumentFileResponse, error) {
	args := ref.Called(id)
	return args.Get(0).([]contracts.DocumentFileResponse), args.Error(1)
}

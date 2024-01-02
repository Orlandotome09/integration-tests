package contactClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockContactClient struct {
	ContactClient
	mock.Mock
}

func (ref *MockContactClient) Get(id string) (*contracts.ContactResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.ContactResponse), args.Error(1)
}

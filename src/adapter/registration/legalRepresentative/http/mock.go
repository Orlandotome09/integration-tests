package legalRepresentativeClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockLegalRepresentativeClient struct {
	LegalRepresentativeClient
	mock.Mock
}

func (ref *MockLegalRepresentativeClient) Get(id string) (*contracts.LegalRepresentativeResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.LegalRepresentativeResponse), args.Error(1)
}

func (ref *MockLegalRepresentativeClient) Search(profileID string) ([]contracts.LegalRepresentativeResponse, error) {
	args := ref.Called(profileID)
	return args.Get(0).([]contracts.LegalRepresentativeResponse), args.Error(1)
}

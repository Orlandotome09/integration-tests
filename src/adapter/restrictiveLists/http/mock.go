package restrictiveListsHttpClient

import "github.com/stretchr/testify/mock"

type MockRestrictiveListsHttpClient struct {
	mock.Mock
}

func (ref *MockRestrictiveListsHttpClient) SearchInternalList(documentFilter string, nameFilter string) (InternalListResponse, error) {
	args := ref.Called(documentFilter, nameFilter)
	return args.Get(0).(InternalListResponse), args.Error(1)
}

func (ref *MockRestrictiveListsHttpClient) SearchPepList(documentNumber string) (*PepResponse, error) {
	args := ref.Called(documentNumber)
	return args.Get(0).(*PepResponse), args.Error(1)
}

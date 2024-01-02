package boardOfDirectorsClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockBoardOfDirectorsClient struct {
	BoardOfDirectorsClient
	mock.Mock
}

func (ref *MockBoardOfDirectorsClient) Search(profileID string) ([]contracts.BoardOfDirectorsResponse, error) {
	args := ref.Called(profileID)
	return args.Get(0).([]contracts.BoardOfDirectorsResponse), args.Error(1)
}

package blacklistClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http/dto"
	"github.com/stretchr/testify/mock"
)

type MockBlacklistClient struct {
	BlackListHttpClient
	mock.Mock
}

func (ref *MockBlacklistClient) Search(documentNumber, partnerId string) (dto.BlacklistResponse, bool, error) {
	args := ref.Called(documentNumber, partnerId)
	return args.Get(0).(dto.BlacklistResponse), args.Bool(1), args.Error(2)
}

package blacklistTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http/dto"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockBlacklistTranslator struct {
	BlacklistTranslator
	mock.Mock
}

func (ref *MockBlacklistTranslator) ToDomain(response *dto.BlacklistResponse) *entity2.BlacklistStatus {
	args := ref.Called(response)
	return args.Get(0).(*entity2.BlacklistStatus)
}

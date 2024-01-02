package accountTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockAccountTranslator struct {
	mock.Mock
}

func (ref *MockAccountTranslator) Translate(accounts []entity.Account) message.Accounts {
	args := ref.Called(accounts)
	return args.Get(0).(message.Accounts)
}

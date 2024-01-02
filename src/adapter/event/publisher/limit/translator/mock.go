package limitMessageTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit/message"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockLimitMessageTranslator struct {
	LimitMessageTranslator
	mock.Mock
}

func (ref *MockLimitMessageTranslator) Translate(profile entity2.Profile, state entity2.State) *message.LimitMessage {
	args := ref.Called(profile)
	return args.Get(0).(*message.LimitMessage)
}

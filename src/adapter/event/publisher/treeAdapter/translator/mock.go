package treeAdapterMessageTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockTreeAdapterMessageTranslator struct {
	TreeAdapterMessageTranslator
	mock.Mock
}

func (ref *MockTreeAdapterMessageTranslator) Translate(profile entity2.Profile, accounts []entity2.Account, addresses []entity2.Address) *message.TreeAdapterMessage {
	args := ref.Called(profile, accounts, addresses)
	return args.Get(0).(*message.TreeAdapterMessage)
}

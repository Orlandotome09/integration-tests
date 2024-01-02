package addressTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockAddressTranslator struct {
	mock.Mock
}

func (ref *MockAddressTranslator) Translate(addresses []entity.Address) message.Addresses {
	args := ref.Called(addresses)
	return args.Get(0).(message.Addresses)
}

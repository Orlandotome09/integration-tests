package registrationNotificationRecipientTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRecipientTranslator struct {
	Translator
	mock.Mock
}

func (ref *MockNotificationRecipientTranslator) Translate(responses []contracts.NotificationRecipientResponse) []entity.NotificationRecipient {
	args := ref.Called(responses)
	return args.Get(0).([]entity.NotificationRecipient)
}

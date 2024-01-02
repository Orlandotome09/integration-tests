package registrationNotificationRecipientTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type Translator interface {
	Translate(responses []contracts.NotificationRecipientResponse) []entity.NotificationRecipient
}

type notificationRecipientTranslator struct{}

func New() Translator {
	return &notificationRecipientTranslator{}
}

func (ref *notificationRecipientTranslator) Translate(responses []contracts.NotificationRecipientResponse) []entity.NotificationRecipient {
	notificationRecipients := make([]entity.NotificationRecipient, len(responses))
	for i, response := range responses {
		notificationType := values.ValidNotificationTypes[response.NotificationType]
		notificationRecipients[i] = entity.NotificationRecipient{
			NotificationType: notificationType,
			EmailTo:          response.EmailTo,
			CopyEmail:        response.CopyEmail,
		}
	}
	return notificationRecipients
}

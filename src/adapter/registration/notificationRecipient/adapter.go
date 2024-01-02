package notificationRecipient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http"
	notificationRecipientTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type notificationRecipientAdapter struct {
	notificationRecipientClient     notificationRecipientClient.NotificationRecipientClient
	notificationRecipientTranslator notificationRecipientTranslator.Translator
}

func NewNotificationRecipientAdapter(notificationRecipientClient notificationRecipientClient.NotificationRecipientClient,
	notificationRecipientTranslator notificationRecipientTranslator.Translator,
) interfaces.NotificationRecipientAdapter {
	return &notificationRecipientAdapter{
		notificationRecipientClient:     notificationRecipientClient,
		notificationRecipientTranslator: notificationRecipientTranslator,
	}
}

func (ref *notificationRecipientAdapter) Search(profileID string) ([]entity.NotificationRecipient, error) {
	responses, err := ref.notificationRecipientClient.Search(profileID)
	if err != nil {
		return nil, err
	}

	notificationRecipients := ref.notificationRecipientTranslator.Translate(responses)

	return notificationRecipients, nil
}

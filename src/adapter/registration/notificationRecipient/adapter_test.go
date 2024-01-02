package notificationRecipient

import (
	notificationRecipientClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/contracts"
	registrationNotificationRecipientesTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	notificationRecipientclient *notificationRecipientClient.MockNotificationRecipientClient
	translator                  *registrationNotificationRecipientesTranslator.MockNotificationRecipientTranslator
	service                     interfaces.NotificationRecipientAdapter
)

func TestMain(m *testing.M) {
	notificationRecipientclient = &notificationRecipientClient.MockNotificationRecipientClient{}
	translator = &registrationNotificationRecipientesTranslator.MockNotificationRecipientTranslator{}
	service = NewNotificationRecipientAdapter(notificationRecipientclient, translator)
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	id := "333"
	responses := []contracts.NotificationRecipientResponse{
		{NotificationType: values.NotificationTypePostWarnings.ToString()},
		{NotificationType: values.NotificationTypeSentOP.ToString()},
	}
	notificationRecipients := []entity.NotificationRecipient{
		{NotificationType: values.NotificationTypePostWarnings},
		{NotificationType: values.NotificationTypeSentOP},
	}

	notificationRecipientclient.On("Search", id).Return(responses, nil)

	translator.On("Translate", responses).Return(notificationRecipients)

	expected := notificationRecipients

	received, err := service.Search(id)

	assert.Equal(t, expected, received)
	assert.Nil(t, err)
}

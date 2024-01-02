package registrationNotificationRecipientTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := New()

	responses := []contracts.NotificationRecipientResponse{
		{
			NotificationType: values.NotificationTypePostWarnings.ToString(),
			EmailTo:          "PostWarningsEmailTo@teste.com.br",
			CopyEmail:        "PostWarningsCopyEmail@teste.com.br",
		},
		{
			NotificationType: values.NotificationTypeSentOP.ToString(),
			EmailTo:          "SentOPEmailTo@teste.com.br",
			CopyEmail:        "SentOPCopyEmail@teste.com.br",
		},
	}

	expected := []entity.NotificationRecipient{
		{
			NotificationType: values.ValidNotificationTypes[responses[0].NotificationType],
			EmailTo:          responses[0].EmailTo,
			CopyEmail:        responses[0].CopyEmail,
		},
		{
			NotificationType: values.ValidNotificationTypes[responses[1].NotificationType],
			EmailTo:          responses[1].EmailTo,
			CopyEmail:        responses[1].CopyEmail,
		},
	}

	received := translator.Translate(responses)

	assert.Equal(t, expected, received)

}

func TestTranslate_no_responses(t *testing.T) {
	translator := New()

	responses := []contracts.NotificationRecipientResponse{}

	expected := []entity.NotificationRecipient{}

	received := translator.Translate(responses)

	assert.Equal(t, expected, received)

}

func TestTranslate_nil_responses(t *testing.T) {
	translator := New()

	var responses []contracts.NotificationRecipientResponse = nil

	expected := []entity.NotificationRecipient{}

	received := translator.Translate(responses)

	assert.Equal(t, expected, received)

}

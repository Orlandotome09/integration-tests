package notificationRecipientConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_When_Tree_Integration_Enabled_Should_Populate_NotificationRecipient(t *testing.T) {
	notificationRecipientService := &mocks.NotificationRecipientService{}
	constructor := notificationRecipientPersonConstructor{notificationRecipientService: notificationRecipientService}

	person := entity.Person{
		EntityID: uuid.New(),
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ProductConfig: &entity.ProductConfig{
				TreeIntegration: true,
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	notificationRecipients := []entity.NotificationRecipient{
		{
			NotificationType: values.NotificationTypePostWarnings,
			EmailTo:          "PostWarningsEmailTo@teste.com.br",
			CopyEmail:        "PostWarningsCopyEmail@teste.com.br",
		},
		{
			NotificationType: values.NotificationTypeSentOP,
			EmailTo:          "SentOPEmailTo@teste.com.br",
			CopyEmail:        "SentOPCopyEmail@teste.com.br",
		},
	}

	notificationRecipientService.On("Search", person.EntityID.String()).Return(notificationRecipients, nil)

	err := constructor.Assemble(personWrapper)

	expected := notificationRecipients

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.NotificationRecipients)
	mock.AssertExpectationsForObjects(t, notificationRecipientService)
}

func Test_Assemble_Should_Not_Populate_NotificationRecipient(t *testing.T) {
	notificationRecipientService := &mocks.NotificationRecipientService{}
	constructor := notificationRecipientPersonConstructor{notificationRecipientService: notificationRecipientService}

	person := entity.Person{
		EntityID:               uuid.New(),
		NotificationRecipients: nil,
	}

	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(personWrapper)

	var expected []entity.NotificationRecipient = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.NotificationRecipients)
	mock.AssertExpectationsForObjects(t, notificationRecipientService)
}

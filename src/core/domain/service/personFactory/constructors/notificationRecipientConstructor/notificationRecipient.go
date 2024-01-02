package notificationRecipientConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type notificationRecipientPersonConstructor struct {
	notificationRecipientService interfaces.NotificationRecipientAdapter
}

func New(service interfaces.NotificationRecipientAdapter) interfaces.PersonConstructor {
	return &notificationRecipientPersonConstructor{notificationRecipientService: service}
}

func (ref *notificationRecipientPersonConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldGetNotificationRecipients() {
		return nil
	}

	notificationRecipients, err := ref.notificationRecipientService.Search(personWrapper.Person.EntityID.String())
	if err != nil {
		return errors.WithStack(err)
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.NotificationRecipients = notificationRecipients

	return nil
}

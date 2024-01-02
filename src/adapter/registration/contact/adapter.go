package contact

import (
	contactClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http"
	contactTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type contactAdapter struct {
	client            contactClient.ContactClient
	contactTranslator contactTranslator.Translator
}

func NewContactAdapter(client contactClient.ContactClient, translator contactTranslator.Translator) interfaces.ContactAdapter {
	return &contactAdapter{
		client:            client,
		contactTranslator: translator,
	}
}

func (ref *contactAdapter) Get(id uuid.UUID) (*entity.Contact, error) {
	resp, err := ref.client.Get(id.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, nil
	}

	contact := &entity.Contact{ProfileID: &resp.ProfileID}

	return contact, nil
}

func (ref *contactAdapter) Search(profileID string) ([]entity.Contact, error) {
	responses, err := ref.client.Search(profileID)
	if err != nil {
		return nil, err
	}

	contacts := ref.contactTranslator.Translate(responses)

	return contacts, nil
}

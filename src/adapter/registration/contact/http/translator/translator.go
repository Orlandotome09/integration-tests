package translator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type Translator interface {
	Translate(responses []contracts.ContactResponse) []entity.Contact
}

type contactTranslator struct{}

func New() Translator {
	return &contactTranslator{}
}

func (ref *contactTranslator) Translate(responses []contracts.ContactResponse) []entity.Contact {
	var contacts []entity.Contact
	for i, response := range responses {
		contact := entity.Contact{
			ContactID:      &responses[i].ContactID,
			ProfileID:      &responses[i].ProfileID,
			Name:           response.Name,
			Email:          response.Email,
			Phone:          response.Phone,
			Phones:         response.Phones.ToDomain(),
			Nationality:    response.Nationality,
			DocumentNumber: response.DocumentNumber,
		}
		contacts = append(contacts, contact)
	}

	return contacts
}

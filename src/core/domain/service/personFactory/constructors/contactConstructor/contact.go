package contactConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type contactPersonConstructor struct {
	contactAdapter interfaces.ContactAdapter
}

func New(adapter interfaces.ContactAdapter) interfaces.PersonConstructor {
	return &contactPersonConstructor{contactAdapter: adapter}
}

func (ref *contactPersonConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldGetContacts() {
		return nil
	}

	contacts, err := ref.contactAdapter.Search(personWrapper.Person.EntityID.String())
	if err != nil {
		return errors.WithStack(err)
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.Contacts = contacts

	return nil
}

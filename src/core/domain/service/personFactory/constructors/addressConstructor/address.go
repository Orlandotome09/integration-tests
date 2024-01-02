package addressConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type addressPersonConstructor struct {
	addressAdapter interfaces.AddressAdapter
}

func New(addressAdapter interfaces.AddressAdapter) interfaces.PersonConstructor {
	return &addressPersonConstructor{addressAdapter: addressAdapter}
}

func (ref *addressPersonConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldGetAddresses() {
		return nil
	}

	addresses, err := ref.addressAdapter.Search(personWrapper.Person.EntityID.String())
	if err != nil {
		return errors.WithStack(err)
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.Addresses = addresses

	return nil
}

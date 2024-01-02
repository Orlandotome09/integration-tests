package enrichmentConstructor

import (
	"bitbucket.org/bexstech/temis-compliance/src/core"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type enrichmentPersonConstructor struct {
	enrichmentAdapter interfaces.EnricherAdapter
}

func New(service interfaces.EnricherAdapter) interfaces.PersonConstructor {
	return &enrichmentPersonConstructor{enrichmentAdapter: service}
}

func (ref *enrichmentPersonConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldValidateCAF() {
		return nil
	}

	document := core.NormalizeDocument(personWrapper.Person.DocumentNumber)

	enrichedInformation, err := ref.enrichmentAdapter.GetEnrichedPerson(document,
		personWrapper.Person.ProfileID.String(),
		personWrapper.Person.PersonType,
		personWrapper.Person.OfferType,
		personWrapper.Person.PartnerID,
		personWrapper.Person.RoleType)
	if err != nil {
		return errors.WithStack(err)
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()

	personWrapper.Person.EnrichedInformation = merge(personWrapper.Person.EnrichedInformation, enrichedInformation)

	return nil
}

// TODO: When Bureau Constructor is replaced by Enrichment, this wont be necessary
func merge(current, new *entity.EnrichedInformation) *entity.EnrichedInformation {

	if current == nil {
		return new
	}

	if new == nil {
		return current
	}

	if len(new.Providers) > 0 {
		current.Providers = new.Providers
	}

	return current
}

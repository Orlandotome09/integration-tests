package pepInformationConstructor

import (
	"bitbucket.org/bexstech/temis-compliance/src/core"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type pepInformationConstructor struct {
	restricitiveListsAdapter interfaces.RestrictiveListsAdapter
}

func New(restricitiveListsAdapter interfaces.RestrictiveListsAdapter) interfaces.PersonConstructor {
	return &pepInformationConstructor{
		restricitiveListsAdapter: restricitiveListsAdapter,
	}
}

func (ref *pepInformationConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldValidatePEP() {
		return nil
	}

	documentNumber := core.NormalizeDocument(personWrapper.Person.DocumentNumber)

	occurrence, err := ref.restricitiveListsAdapter.OccurrenceInPepList(documentNumber)
	if err != nil {
		return nil
	}

	if occurrence == nil {
		return nil
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.PEPInformation = occurrence

	return nil
}

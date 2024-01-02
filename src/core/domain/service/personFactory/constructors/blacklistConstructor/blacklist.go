package blacklistConstructor

import (
	"bitbucket.org/bexstech/temis-compliance/src/core"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type blacklistPersonConstructor struct {
	restricitiveListsAdapter interfaces.RestrictiveListsAdapter
}

func New(restricitiveListsAdapter interfaces.RestrictiveListsAdapter) interfaces.PersonConstructor {
	return &blacklistPersonConstructor{
		restricitiveListsAdapter: restricitiveListsAdapter,
	}
}

func (ref *blacklistPersonConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldValidateBlacklist() {
		return nil
	}

	documentNumber := core.NormalizeDocument(personWrapper.Person.DocumentNumber)
	name := personWrapper.Person.Name

	occurrence, err := ref.restricitiveListsAdapter.OccurrenceInBlackList(documentNumber, name)
	if err != nil {
		return nil
	}

	if occurrence == nil {
		return nil
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.BlacklistStatus = occurrence

	return nil
}

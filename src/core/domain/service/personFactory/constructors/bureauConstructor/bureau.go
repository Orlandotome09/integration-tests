package bureauConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
	"time"
)

type bureauPersonConstructor struct {
	bureauService interfaces.BureauService
}

func New(service interfaces.BureauService) interfaces.PersonConstructor {
	return &bureauPersonConstructor{bureauService: service}
}

func (ref *bureauPersonConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldGetBureauInformation() {
		return nil
	}

	bureauResponse, err := ref.bureauService.GetBureauStatus(personWrapper.Person)
	if err != nil {
		return errors.WithStack(err)
	}

	if bureauResponse != nil && bureauResponse.BirthDate != "" {
		dateOfBirth, err := time.Parse("02/01/2006", bureauResponse.BirthDate)
		if err != nil {
			return errors.WithStack(err)
		}
		if personWrapper.Person.Individual == nil {
			personWrapper.Person.Individual = &entity.Individual{}
		}
		personWrapper.Person.Individual.DateOfBirth = &dateOfBirth
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.EnrichedInformation = bureauResponse

	return nil
}

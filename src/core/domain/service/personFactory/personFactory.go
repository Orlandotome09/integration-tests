package personFactory

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type personFactory struct {
	overrideService                interfaces.OverrideService
	cadastralValidationConstructor interfaces.PersonConstructor
	personRulesConstructor         interfaces.PersonConstructor
	bureauConstructor              interfaces.PersonConstructor
	personConstructors             []interfaces.PersonConstructor
}

func New(overrideService interfaces.OverrideService,
	cadastralValidationConstructor interfaces.PersonConstructor,
	personRulesConstructor interfaces.PersonConstructor,
	bureauConstructor interfaces.PersonConstructor,
	constructors []interfaces.PersonConstructor) interfaces.PersonFactory {

	return &personFactory{
		overrideService:                overrideService,
		cadastralValidationConstructor: cadastralValidationConstructor,
		personRulesConstructor:         personRulesConstructor,
		bureauConstructor:              bureauConstructor,
		personConstructors:             constructors,
	}
}

func (ref *personFactory) Build(person entity.Person) (*entity.Person, error) {

	personWrapper := entity.PersonWrapper{
		Person: person,
	}

	err := ref.cadastralValidationConstructor.Assemble(&personWrapper)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = ref.bureauConstructor.Assemble(&personWrapper)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	constructorsGroup := new(errgroup.Group)
	for _, constructor := range ref.personConstructors {
		auxConstructor := constructor
		constructorsGroup.Go(func() error {
			return auxConstructor.Assemble(&personWrapper)
		})
	}
	if err := constructorsGroup.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	err = ref.personRulesConstructor.Assemble(&personWrapper)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	overrides, err := ref.overrideService.FindByEntityID(personWrapper.Person.EntityID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	personWrapper.Person.Overrides = overrides

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &personWrapper.Person, nil
}

package profileFactory

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type profileFactory struct {
	profileAdapter      interfaces.ProfileAdapter
	personFactory       interfaces.PersonFactory
	rulesConstructor    interfaces.ProfileConstructor
	profileConstructors []interfaces.ProfileConstructor
}

func New(
	profileAdapter interfaces.ProfileAdapter,
	personFactory interfaces.PersonFactory,
	rulesConstructor interfaces.ProfileConstructor,
	profileConstructors []interfaces.ProfileConstructor) interfaces.ProfileFactory {
	return &profileFactory{
		personFactory:       personFactory,
		profileAdapter:      profileAdapter,
		rulesConstructor:    rulesConstructor,
		profileConstructors: profileConstructors,
	}
}

func (ref *profileFactory) Build(profileID uuid.UUID) (*entity.Profile, error) {

	profile, err := ref.profileAdapter.Get(profileID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if profile == nil {
		return nil, nil
	}

	person, err := ref.personFactory.Build(profile.Person)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if person != nil {
		profile.Person = *person
	}

	profileWrapper := entity.ProfileWrapper{
		Profile: *profile,
	}

	constructorsGroup := new(errgroup.Group)
	for _, constructor := range ref.profileConstructors {
		auxConstructor := constructor
		constructorsGroup.Go(func() error {
			return auxConstructor.Assemble(&profileWrapper)
		})
	}
	if err := constructorsGroup.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	err = ref.rulesConstructor.Assemble(&profileWrapper)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &profileWrapper.Profile, nil

}

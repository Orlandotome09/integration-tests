package legalRepresentativesConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type legalRepresentativesConstructor struct {
	legalRepresentativeAdapter interfaces.LegalRepresentativeAdapter
}

func New(legalRepresentativeAdapter interfaces.LegalRepresentativeAdapter) interfaces.ProfileConstructor {
	return &legalRepresentativesConstructor{legalRepresentativeAdapter: legalRepresentativeAdapter}
}

func (ref *legalRepresentativesConstructor) Assemble(profileWrapper *entity.ProfileWrapper) error {

	if !profileWrapper.Profile.ShouldGetLegalRepresentatives() {
		return nil
	}

	legalRepresentatives, err := ref.legalRepresentativeAdapter.Search(*profileWrapper.Profile.ProfileID)
	if err != nil {
		return errors.WithStack(err)
	}

	profileWrapper.Mutex.Lock()
	defer profileWrapper.Mutex.Unlock()
	profileWrapper.Profile.LegalRepresentatives = legalRepresentatives

	return nil

}

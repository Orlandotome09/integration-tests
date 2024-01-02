package ownershipStructureConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type ownershipStructureConstructor struct {
	ownershipStructureService interfaces.OwnershipStructureService
}

func New(ownershipStructureService interfaces.OwnershipStructureService) interfaces.ProfileConstructor {
	return &ownershipStructureConstructor{ownershipStructureService: ownershipStructureService}
}

func (ref *ownershipStructureConstructor) Assemble(profileWrapper *entity.ProfileWrapper) error {

	if !profileWrapper.Profile.ShouldGetOwnershipStructure() {
		return nil
	}

	ownershipStructure, err := ref.ownershipStructureService.GetManuallyFilled(profileWrapper.Profile.ProfileID.String())
	if err != nil {
		return errors.WithStack(err)
	}

	enriched, err := ref.ownershipStructureService.GetEnriched(profileWrapper.Profile.DocumentNumber, profileWrapper.Profile.OfferType, profileWrapper.Profile.PartnerID)
	if err != nil {
		return errors.WithStack(err)
	}

	profileWrapper.Mutex.Lock()
	defer profileWrapper.Mutex.Unlock()
	profileWrapper.Profile.OwnershipStructure = ownershipStructure

	if profileWrapper.Profile.EnrichedInformation == nil {
		profileWrapper.Profile.EnrichedInformation = &entity.EnrichedInformation{}
	}

	profileWrapper.Profile.EnrichedInformation.OwnershipStructure = enriched

	return nil

}

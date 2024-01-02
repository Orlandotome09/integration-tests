package ownershipStructure

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ownershipStructureService struct {
	adapterEnrichment   interfaces.OwnershipStructureAdapter
	adapterRegistration interfaces.OwnershipStructureAdapter
}

func New(adapterEnrichment interfaces.OwnershipStructureAdapter,
	adapterRegistration interfaces.OwnershipStructureAdapter) interfaces.OwnershipStructureService {
	return &ownershipStructureService{
		adapterEnrichment:   adapterEnrichment,
		adapterRegistration: adapterRegistration,
	}
}

func (ref *ownershipStructureService) GetEnriched(legalEntityID, offerType, partnerID string) (*entity.OwnershipStructure, error) {
	return ref.adapterEnrichment.Get(legalEntityID, offerType, partnerID)
}

func (ref *ownershipStructureService) GetManuallyFilled(profileID string) (*entity.OwnershipStructure, error) {
	return ref.adapterRegistration.Get(profileID, "", "")
}

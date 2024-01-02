package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type OwnershipStructureService interface {
	GetEnriched(legalEntityID, offerType, partnerID string) (*entity.OwnershipStructure, error)
	GetManuallyFilled(profileID string) (*entity.OwnershipStructure, error)
}

type OwnershipStructureAdapter interface {
	Get(id, offerType, partnerID string) (*entity.OwnershipStructure, error)
}

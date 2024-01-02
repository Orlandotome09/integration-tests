package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type EnricherAdapter interface {
	GetEnrichedPerson(documentNumber, profileID, personType, offerType, partnerID, roleTye string) (*entity.EnrichedInformation, error)
}

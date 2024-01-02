package _interfacesEnrichment

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type EnrichedInformationTranslator interface {
	Translate(documentNumber string, profileID uuid.UUID, response []byte) (*entity.EnrichedInformation, error)
}

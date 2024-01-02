package enrichedIndividualTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/_interfacesEnrichment"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/individual/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type enrichedIndividualTranslator struct{}

func New() _interfacesEnrichment.EnrichedInformationTranslator {
	return &enrichedIndividualTranslator{}
}

func (ref *enrichedIndividualTranslator) Translate(documentNumber string, profileID uuid.UUID, response []byte) (*entity.EnrichedInformation, error) {

	var individualResponse contracts.IndividualResponse
	err := json.Unmarshal(response, &individualResponse)
	if err != nil {
		return nil, errors.Errorf("could not convert, response is not an individual: %+v", response)
	}

	return &entity.EnrichedInformation{
		BureauStatus: ref.translateSituationToStatus(individualResponse.Situation),
		EnrichedIndividual: entity.EnrichedIndividual{
			Name:      individualResponse.Name,
			BirthDate: individualResponse.BirthDate,
		},
	}, nil
}

func (ref *enrichedIndividualTranslator) translateSituationToStatus(situation int) string {
	switch situation {
	case 0:
		return "DECEASED"
	case 1:
		return "REGULAR"
	case 2:
		return "PENDING_REGULARIZATION"
	case 3:
		return "SUSPENDED"
	case 4:
		return "CANCELLED"
	case 5:
		return "NULL"
	default:
		return "UNKNOWN"
	}
}

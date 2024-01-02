package enrichedLegalEntityTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/_interfacesEnrichment"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/legalEntity/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type enrichedLegalEntityTranslator struct{}

func New() _interfacesEnrichment.EnrichedInformationTranslator {
	return &enrichedLegalEntityTranslator{}
}

func (ref *enrichedLegalEntityTranslator) Translate(documentNumber string, profileID uuid.UUID, response []byte) (*entity.EnrichedInformation, error) {

	var legalEntityResponse contracts.LegalEntityResponse
	err := json.Unmarshal(response, &legalEntityResponse)
	if err != nil {
		return nil, errors.Errorf("could not convert, response is not a legal entity: %+v", response)
	}

	directors := make([]entity.Director, 0)

	for _, director := range legalEntityResponse.BoardOfDirectors {

		directorID := uuid.NewSHA1(core.GetUuidNamespace(), []byte(documentNumber+director.DocumentNumber))

		director := entity.Director{
			DirectorID: directorID,
			Role:       director.Role,
			Person: entity.Person{
				EntityID:       directorID,
				ProfileID:      profileID,
				Name:           director.Name,
				DocumentNumber: director.DocumentNumber,
				PersonType:     values.PersonTypeIndividual,
				RoleType:       values.RoleTypeDirector,
				EntityType:     values.EntityTypeDirector,
			},
		}

		directors = append(directors, director)
	}

	return &entity.EnrichedInformation{
		BureauStatus: ref.translateSituationToStatus(legalEntityResponse.Situation),
		EnrichedCompany: entity.EnrichedCompany{
			LegalName:        legalEntityResponse.LegalName,
			EconomicActivity: legalEntityResponse.CNAE,
			OpeningDate:      legalEntityResponse.OpeningDate,
			LegalNature:      legalEntityResponse.LegalNature,
			BoardOfDirectors: directors,
		},
	}, nil
}

func (ref *enrichedLegalEntityTranslator) translateSituationToStatus(situation int) string {
	switch situation {
	case 2:
		return "REGULAR"
	default:
		return "NOT_REGULAR"
	}
}

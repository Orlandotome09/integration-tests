package enricherTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type EnricherTranslator interface {
	Translate(response enricherContracts.EnricherResponse) *entity.EnrichedInformation
}

type enricherTranslator struct{}

func New() EnricherTranslator {
	return &enricherTranslator{}
}

func (ref *enricherTranslator) Translate(response enricherContracts.EnricherResponse) *entity.EnrichedInformation {

	bureauStatus := ""
	providers := make([]entity.Provider, 0)

	for _, providerResponse := range response.Providers {
		switch providerResponse.ProviderName {
		case values.IndividualBureauEnricher.String():
			bureauStatus = ref.translateIndividualSituationToStatus(response.Individual.Situation)
		case values.LegalEntityBureauEnricher.String():
			bureauStatus = ref.translateCompanySituationToStatus(response.Company.Situation)
		}

		provider := entity.Provider{
			ProviderName: providerResponse.ProviderName,
			RequestID:    providerResponse.RequestID,
			Status:       providerResponse.Status,
		}
		providers = append(providers, provider)
	}

	return &entity.EnrichedInformation{
		BureauStatus:       bureauStatus,
		EnrichedIndividual: ref.translateIndividual(response),
		EnrichedCompany:    ref.translateCompany(response),
		Providers:          providers,
	}
}

func (ref *enricherTranslator) translateIndividualSituationToStatus(situation int) string {
	switch situation {
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

func (ref *enricherTranslator) translateCompanySituationToStatus(situation int) string {
	switch situation {
	case 2:
		return "REGULAR"
	default:
		return "NOT_REGULAR"
	}
}

func (ref *enricherTranslator) translateIndividual(response enricherContracts.EnricherResponse) entity.EnrichedIndividual {
	if response.Person.Type != values.PersonTypeIndividual {
		return entity.EnrichedIndividual{}
	}

	return entity.EnrichedIndividual{
		Name:      response.Person.Name,
		BirthDate: response.Individual.BirthDate,
	}
}

func (ref *enricherTranslator) translateCompany(response enricherContracts.EnricherResponse) entity.EnrichedCompany {
	if response.Person.Type != values.PersonTypeCompany {
		return entity.EnrichedCompany{}
	}

	directors := make([]entity.Director, 0)

	for _, director := range response.Company.BoardOfDirectors {

		directorID := uuid.NewSHA1(core.GetUuidNamespace(), []byte(response.DocumentNumber+director.DocumentNumber))

		director := entity.Director{
			DirectorID: directorID,
			Role:       director.Role,
			Person: entity.Person{
				EntityID:       directorID,
				Name:           director.Name,
				DocumentNumber: director.DocumentNumber,
				PersonType:     values.PersonTypeIndividual,
				RoleType:       values.RoleTypeDirector,
				EntityType:     values.EntityTypeDirector,
			},
		}

		directors = append(directors, director)
	}

	return entity.EnrichedCompany{
		LegalName:        response.Name,
		EconomicActivity: response.Company.CNAE,
		OpeningDate:      response.Company.OpeningDate,
		LegalNature:      response.Company.LegalNature,
		BoardOfDirectors: directors,
	}
}

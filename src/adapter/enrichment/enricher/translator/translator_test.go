package enricherTranslator

import (
	enricherContracts "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslateIndividual(t *testing.T) {
	translator := New()
	profileID := uuid.New()

	responseHttpClient := enricherContracts.EnricherResponse{
		Person: enricherContracts.Person{
			EntityID:       profileID,
			Role:           values.RoleTypeCustomer,
			Type:           values.PersonTypeIndividual,
			Name:           "JOAO DA SILVA",
			DocumentNumber: "58349083069",
			Individual: &enricherContracts.IndividualResponse{
				Nationality: "BRASILEIRO",
				BirthDate:   "30/11/1980",
				Situation:   1,
			},
		},
		Providers: []enricherContracts.Provider{{
			ProviderName: values.IndividualBureauEnricher.String(),
		}},
	}

	expected := &entity.EnrichedInformation{
		BureauStatus: "REGULAR",
		EnrichedIndividual: entity.EnrichedIndividual{
			Name:      responseHttpClient.Name,
			BirthDate: responseHttpClient.Individual.BirthDate,
		},
		Providers: []entity.Provider{{
			ProviderName: values.IndividualBureauEnricher.String(),
		}},
	}

	received := translator.Translate(responseHttpClient)

	assert.Equal(t, expected, received)
}

func TestTranslateCompany(t *testing.T) {
	translator := New()
	profileID := uuid.New()

	responseHttpClient := enricherContracts.EnricherResponse{
		Person: enricherContracts.Person{
			EntityID:       profileID,
			Role:           values.RoleTypeMerchant,
			Type:           values.PersonTypeCompany,
			Name:           "EMPRESA TESTE LTDA",
			DocumentNumber: "45614677000158",
			Company: &enricherContracts.LegalEntityResponse{
				BusinessName: "NOME FANTASIA",
				CNAE:         "1531-9/02",
				OpeningDate:  "01/01/1980",
				LegalNature:  "MEI",
				Situation:    2,
			},
		},
		Providers: []enricherContracts.Provider{{
			ProviderName: values.LegalEntityBureauEnricher.String(),
		}},
	}

	expected := &entity.EnrichedInformation{
		BureauStatus: "REGULAR",

		EnrichedCompany: entity.EnrichedCompany{
			LegalName:        responseHttpClient.Person.Name,
			EconomicActivity: responseHttpClient.Company.CNAE,
			OpeningDate:      responseHttpClient.Company.OpeningDate,
			LegalNature:      responseHttpClient.Company.LegalNature,
			BoardOfDirectors: []entity.Director{},
		},
		Providers: []entity.Provider{{
			ProviderName: values.LegalEntityBureauEnricher.String(),
		}},
	}

	received := translator.Translate(responseHttpClient)

	assert.Equal(t, expected, received)
}

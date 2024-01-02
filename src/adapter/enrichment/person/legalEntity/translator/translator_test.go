package enrichedLegalEntityTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/legalEntity/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := New()

	response := &contracts.LegalEntityResponse{
		LegalName:   "Tânia e Lorenzo Eletrônica Ltda",
		OpeningDate: "210/02/2016",
		Situation:   2,
		LegalNature: "2111",
		CNAE:        "SomeCNAE",
		BoardOfDirectors: []contracts.Director{
			{
				DocumentNumber: "123",
				Name:           "Some Director",
				Role:           "10",
				PersonType:     "INDIVIDUAL",
			},
		},
	}

	responseBytes := new(bytes.Buffer)
	json.NewEncoder(responseBytes).Encode(response)

	profileID := uuid.New()
	expected := &entity.EnrichedInformation{
		BureauStatus: "REGULAR",
		EnrichedCompany: entity.EnrichedCompany{
			LegalName:        response.LegalName,
			EconomicActivity: "SomeCNAE",
			OpeningDate:      response.OpeningDate,
			LegalNature:      "2111",
			BoardOfDirectors: []entity.Director{
				{
					Role: "10",
					Person: entity.Person{
						DocumentNumber: "123",
						Name:           "Some Director",
						PersonType:     values.PersonTypeIndividual,
						RoleType:       values.RoleTypeDirector,
						EntityType:     values.EntityTypeDirector,
						ProfileID:      profileID,
					},
				},
			},
		},
	}

	received, err := translator.Translate("123", profileID, responseBytes.Bytes())

	expected.EnrichedCompany.BoardOfDirectors[0].DirectorID = received.EnrichedCompany.BoardOfDirectors[0].DirectorID
	expected.EnrichedCompany.BoardOfDirectors[0].Person.EntityID = received.EnrichedCompany.BoardOfDirectors[0].DirectorID

	assert.Nil(t, err)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

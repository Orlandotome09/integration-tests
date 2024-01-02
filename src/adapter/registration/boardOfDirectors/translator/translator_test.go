package boardOfDirectorsTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {

	translator := New()
	birthDate, _ := time.Parse("2006-01-02", "2006-01-02")
	pep := true
	response := contracts.BoardOfDirectorsResponse{
		DirectorID:     uuid.New(),
		ProfileID:      uuid.New(),
		FullName:       "COMPLETE NAME",
		DocumentNumber: "123",
		Nationality:    "BR",
		DateOfBirth:    "2006-01-02",
		Pep:            &pep,
	}

	expected := &entity.Director{
		DirectorID: response.DirectorID,
		Person: entity.Person{
			DocumentNumber: "123",
			Name:           "COMPLETE NAME",
			PersonType:     "INDIVIDUAL",
			ProfileID:      response.ProfileID,
			EntityID:       response.DirectorID,
			EntityType:     "DIRECTOR",
			RoleType:       "DIRECTOR",
			Individual: &entity.Individual{
				Nationality: "BR",
				DateOfBirth: &birthDate,
				Pep:         &pep,
			},
		},
	}

	received, err := translator.Translate(response)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	assert.Nil(t, err)
}

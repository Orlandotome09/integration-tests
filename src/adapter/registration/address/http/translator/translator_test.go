package registrationAddressesTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {
	translator := New()

	responses := []contracts.AddressResponse{
		{
			AddressID:    uuid.New(),
			ProfileID:    uuid.New(),
			Type:         "X",
			ZipCode:      "333",
			Street:       "Rua Antiguidades",
			Number:       "333",
			Complement:   "bloco C",
			Neighborhood: "Vila Yara",
			City:         "São Paulo",
			StateCode:    "SP",
			CountryCode:  "BRA",
		},
		{
			AddressID:    uuid.New(),
			ProfileID:    uuid.New(),
			Type:         "Y",
			ZipCode:      "444",
			Street:       "Rua dos Antiquários",
			Number:       "444",
			Complement:   "bloco D",
			Neighborhood: "Vila Suzi",
			City:         "Rio de Janeiro",
			StateCode:    "RJ",
			CountryCode:  "BRA",
		},
	}

	expected := []entity.Address{
		{
			AddressID:    &responses[0].AddressID,
			ProfileID:    &responses[0].ProfileID,
			Type:         responses[0].Type,
			ZipCode:      responses[0].ZipCode,
			Street:       responses[0].Street,
			Number:       responses[0].Number,
			Complement:   responses[0].Complement,
			Neighborhood: responses[0].Neighborhood,
			City:         responses[0].City,
			StateCode:    responses[0].StateCode,
			CountryCode:  responses[0].CountryCode,
		},
		{
			AddressID:    &responses[1].AddressID,
			ProfileID:    &responses[1].ProfileID,
			Type:         responses[1].Type,
			ZipCode:      responses[1].ZipCode,
			Street:       responses[1].Street,
			Number:       responses[1].Number,
			Complement:   responses[1].Complement,
			Neighborhood: responses[1].Neighborhood,
			City:         responses[1].City,
			StateCode:    responses[1].StateCode,
			CountryCode:  responses[1].CountryCode,
		},
	}

	received := translator.Translate(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslate_no_responses(t *testing.T) {
	translator := New()

	responses := []contracts.AddressResponse{}

	expected := []entity.Address{}

	received := translator.Translate(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslate_nil_responses(t *testing.T) {
	translator := New()

	var responses []contracts.AddressResponse = nil

	expected := []entity.Address{}

	received := translator.Translate(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestSortByUpdateDateDesc(t *testing.T) {

	responses := []contracts.AddressResponse{
		{
			Street:    "Rua Helena",
			UpdatedAt: time.Date(2020, 2, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			Street:    "Rua dos Antiquários",
			UpdatedAt: time.Date(2010, 2, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			Street:    "Rua Augusta",
			UpdatedAt: time.Date(2015, 2, 21, 0, 0, 0, 0, time.UTC),
		},
	}

	expected := []contracts.AddressResponse{responses[0], responses[2], responses[1]}

	received := sortByUpdateDateDesc(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

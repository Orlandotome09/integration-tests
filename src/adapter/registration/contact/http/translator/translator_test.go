package translator

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"reflect"
	"testing"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/contracts"
	"github.com/google/uuid"
)

func TestTranslate(t *testing.T) {
	translator := New()

	responses := []contracts.ContactResponse{
		{
			ContactID:      uuid.New(),
			ProfileID:      uuid.New(),
			Name:           "Residential",
			DocumentNumber: "111",
			Email:          "home@mail.com",
			Phone:          "01123456789",
			Nationality:    "BR",
			Phones:         contracts.Phones{},
		},
		{
			ContactID:      uuid.New(),
			ProfileID:      uuid.New(),
			Name:           "Work",
			DocumentNumber: "222",
			Email:          "work@mail.com",
			Phone:          "09923456789",
			Nationality:    "BR",
			Phones:         contracts.Phones{},
		},
	}

	expected := []entity2.Contact{
		{
			ContactID:      &responses[0].ContactID,
			ProfileID:      &responses[0].ProfileID,
			Name:           responses[0].Name,
			Email:          responses[0].Email,
			Phone:          responses[0].Phone,
			Nationality:    responses[0].Nationality,
			DocumentNumber: responses[0].DocumentNumber,
			Phones:         []entity2.Phone{},
		},
		{
			ContactID:      &responses[1].ContactID,
			ProfileID:      &responses[1].ProfileID,
			Name:           responses[1].Name,
			Email:          responses[1].Email,
			Phone:          responses[1].Phone,
			Nationality:    responses[1].Nationality,
			DocumentNumber: responses[1].DocumentNumber,
			Phones:         []entity2.Phone{},
		},
	}

	received := translator.Translate(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslate_no_responses(t *testing.T) {
	translator := New()

	var responses []contracts.ContactResponse

	var expected []entity2.Contact

	received := translator.Translate(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslate_nil_responses(t *testing.T) {
	translator := New()

	var responses []contracts.ContactResponse = nil

	var expected []entity2.Contact

	received := translator.Translate(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

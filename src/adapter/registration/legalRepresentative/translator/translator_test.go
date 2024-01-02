package legalRepresentativeTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {

	translator := New()
	legalRepresentativeID := uuid.New()
	profileID := uuid.New()
	partnerID := uuid.New()
	birthDate, _ := time.Parse("02/01/2006", "01/01/2000")
	pep := true
	response := contracts.LegalRepresentativeResponse{
		LegalRepresentativeID: legalRepresentativeID.String(),
		ProfileID:             profileID.String(),
		PartnerID:             partnerID.String(),
		FullName:              "COMPLETE NAME",
		DocumentNumber:        "123",
		Email:                 "some@email.com",
		Phone:                 "123456",
		Nationality:           "BR",
		BirthDate:             "01/01/2000",
		OfferType:             "OFFER",
		Pep:                   &pep,
	}

	expected := &entity.LegalRepresentative{
		LegalRepresentativeID: legalRepresentativeID,
		Person: entity.Person{
			DocumentNumber: "123",
			Name:           "COMPLETE NAME",
			PersonType:     "INDIVIDUAL",
			Email:          "some@email.com",
			PartnerID:      partnerID.String(),
			OfferType:      "OFFER",
			ProfileID:      profileID,
			EntityID:       legalRepresentativeID,
			EntityType:     "LEGAL_REPRESENTATIVE",
			RoleType:       "LEGAL_REPRESENTATIVE",
			Individual: &entity.Individual{
				Nationality: "BR",
				DateOfBirth: &birthDate,
				Pep:         &pep,
				Phones: []entity.Phone{
					{
						Number: "123456",
					},
				},
			},
		},
	}

	received, err := translator.Translate(response)

	assert.Nil(t, err)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

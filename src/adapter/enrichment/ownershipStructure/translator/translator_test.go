package translator

import (
	enrichedOwnershipStructureHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {
	translator := New()
	birthDate, _ := time.Parse("02/01/2006", "01/01/2000")
	legalEntityID := uuid.New()
	shareholder1ID := "111"
	uuid1 := uuid.NewSHA1(core.GetUuidNamespace(), []byte(shareholder1ID))
	shareholder2ID := "222"
	uuid2 := uuid.NewSHA1(core.GetUuidNamespace(), []byte(shareholder2ID))
	shareholder3ID := "333"
	uuid3 := uuid.NewSHA1(core.GetUuidNamespace(), []byte(shareholder3ID))

	response := enrichedOwnershipStructureHttpClient.OwnershipStructureResponse{
		LegalEntityID:             legalEntityID.String(),
		FinalBeneficiariesCounted: 5,
		ShareholdingSum:           100,
		Shareholders: enrichedOwnershipStructureHttpClient.Shareholders{
			enrichedOwnershipStructureHttpClient.Shareholder{
				ParentLegalEntity: "",
				Shareholding:      10,
				Role:              "",
				Type:              "LEGAL_ENTITY",
				Name:              "COCACOLA",
				DocumentNumber:    shareholder1ID,
				Nationality:       "BRA",
				BirthDate:         "01/01/2000",
				Shareholders: enrichedOwnershipStructureHttpClient.Shareholders{
					enrichedOwnershipStructureHttpClient.Shareholder{
						ParentLegalEntity: shareholder1ID,
						Shareholding:      100,
						Role:              "",
						Type:              "INDIVIDUAL",
						Name:              "JOSE",
						DocumentNumber:    shareholder2ID,
						Nationality:       "BRA",
						BirthDate:         "01/01/2000",
					},
				},
			},
			enrichedOwnershipStructureHttpClient.Shareholder{
				ParentLegalEntity: "",
				Shareholding:      90,
				Role:              "",
				Type:              "INDIVIDUAL",
				Name:              "LUIS",
				DocumentNumber:    shareholder3ID,
				Nationality:       "BRA",
				BirthDate:         "01/01/2000",
			},
		},
	}

	expected := &entity.OwnershipStructure{
		FinalBeneficiariesCount: 5,
		ShareholdingSum:         100,
		Shareholders: []entity.Shareholder{
			{
				ShareholderID:    &uuid1,
				OwnershipPercent: 10,
				Person: entity.Person{
					DocumentNumber: shareholder1ID,
					Name:           "COCACOLA",
					PersonType:     "COMPANY",
					EntityID:       uuid1,
					EntityType:     "SHAREHOLDER",
					RoleType:       "SHAREHOLDER",
					Company: &entity.Company{
						LegalName: "COCACOLA",
					},
				},
			},
			{
				ShareholderID:    &uuid2,
				OwnershipPercent: 100,
				Person: entity.Person{
					DocumentNumber: shareholder2ID,
					Name:           "JOSE",
					PersonType:     "INDIVIDUAL",
					EntityID:       uuid2,
					EntityType:     "SHAREHOLDER",
					RoleType:       "SHAREHOLDER",
					Individual: &entity.Individual{
						DateOfBirth:         &birthDate,
						DateOfBirthInputted: &birthDate,
					},
				},
			},
			{
				ShareholderID:    &uuid3,
				OwnershipPercent: 90,
				Person: entity.Person{
					DocumentNumber: shareholder3ID,
					Name:           "LUIS",
					PersonType:     "INDIVIDUAL",
					EntityID:       uuid3,
					EntityType:     "SHAREHOLDER",
					RoleType:       "SHAREHOLDER",
					Individual: &entity.Individual{
						DateOfBirth:         &birthDate,
						DateOfBirthInputted: &birthDate,
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

func TestRemoveDuplicated(t *testing.T) {
	shareholders := []entity.Shareholder{
		{Person: entity.Person{DocumentNumber: "1"}},
		{Person: entity.Person{DocumentNumber: "2"}},
		{Person: entity.Person{DocumentNumber: "2"}},
	}

	expected := []entity.Shareholder{
		{Person: entity.Person{DocumentNumber: "1"}},
		{Person: entity.Person{DocumentNumber: "2"}},
	}

	received := removeDuplicated(shareholders)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

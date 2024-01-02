package treeAdapterMessageTranslator

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message/types"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/accountTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/addressTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/companyTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/individualTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/timeTranslator"
	"github.com/google/uuid"
)

func TestTranslate_individual(t *testing.T) {
	timeTranslatorInstance := &timeTranslator.MockTimeTranslator{}
	accountTranslatorInstance := &accountTranslator.MockAccountTranslator{}
	addressTranslatorInstance := &addressTranslator.MockAddressTranslator{}
	individualTranslatorInstance := &individualTranslator.MockIndividualTranslator{}
	companyTranslatorInstance := &companyTranslator.MockCompanyTranslator{}
	translator := New(timeTranslatorInstance, accountTranslatorInstance, addressTranslatorInstance, individualTranslatorInstance, companyTranslatorInstance)

	profileID := uuid.New()
	createdDate := time.Date(2020, 0o1, 0o1, 1, 1, 1, 1, time.UTC)
	dateOfBirthInputed := time.Date(2021, 9, 9, 1, 1, 1, 1, time.UTC)
	pep := false

	profile := entity2.Profile{
		ProfileID: &profileID,
		Person: entity2.Person{
			DocumentNumber: "12345678911",
			PersonType:     values.PersonTypeIndividual,
			PartnerID:      "Partner",
			Individual: &entity2.Individual{
				FirstName:           "Patricia",
				LastName:            "Lemos",
				DateOfBirthInputted: &dateOfBirthInputed,
				Nationality:         "BRA",
				Pep:                 &pep,
			},
		},
		CreatedAt: createdDate,
	}

	accounts := []entity2.Account{{AccountNumber: "xxx"}}
	addresses := []entity2.Address{{Street: "Street1"}}

	translatedAccounts := message.Accounts{{AgencyNumber: "123"}, {AgencyNumber: "456"}}
	translatedAddresses := message.Addresses{{Street: "Street1"}}
	translatedName := "Patricia"
	translatedDateOfBirth := "01012021"
	translatedNationality := "BRA"
	translatedIncome := 10000.0
	translatedAssets := 10000.0
	date := time.Now().UTC()

	timeTranslatorInstance.On("TranslateTime", profile.CreatedAt).Return("01012020")
	timeTranslatorInstance.On("GenerateMessageDate").Return(date)
	accountTranslatorInstance.On("Translate", accounts).Return(translatedAccounts)
	addressTranslatorInstance.On("Translate", addresses).Return(translatedAddresses)
	individualTranslatorInstance.On("TranslateName", profile.Individual).Return(translatedName)
	individualTranslatorInstance.On("TranslateDateOfBirth", profile.Individual).Return(translatedDateOfBirth)
	individualTranslatorInstance.On("TranslateNationality", profile.Individual).Return(translatedNationality)
	individualTranslatorInstance.On("TranslatePep", profile.Individual).Return(false)
	individualTranslatorInstance.On("TranslateIncome", profile.Individual).Return(translatedIncome)
	individualTranslatorInstance.On("TranslateAssets", profile.Individual).Return(translatedAssets)

	expected := &message.TreeAdapterMessage{
		ProfileID:          profile.ProfileID.String(),
		Partner:            profile.PartnerID,
		ProfileCreatedDate: "01012020",
		Negotiation:        types.NegotiationTypeMesa,
		Person: message.Person{
			ID:                     profile.DocumentNumber,
			Type:                   types.PersonTypeIndividual,
			Name:                   translatedName,
			BirthDate:              translatedDateOfBirth,
			Nationality:            translatedNationality,
			DigitalSign:            false,
			Addresses:              translatedAddresses,
			Accounts:               translatedAccounts,
			Contacts:               message.Contacts{},
			NotificationRecipients: message.NotificationRecipients{},
			Individual: message.Individual{
				Pep:    false,
				Income: translatedIncome,
				Assets: translatedAssets,
			},
		},
		Date: date,
	}

	received := translator.Translate(profile, accounts, addresses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslate_company(t *testing.T) {
	timeTranslatorInstance := &timeTranslator.MockTimeTranslator{}
	accountTranslatorInstance := &accountTranslator.MockAccountTranslator{}
	addressTranslatorInstance := &addressTranslator.MockAddressTranslator{}
	individualTranslatorInstance := &individualTranslator.MockIndividualTranslator{}
	companyTranslatorInstance := &companyTranslator.MockCompanyTranslator{}
	translator := New(timeTranslatorInstance, accountTranslatorInstance, addressTranslatorInstance, individualTranslatorInstance, companyTranslatorInstance)

	profileID := uuid.New()
	createdDate := time.Date(2020, 0o1, 0o1, 1, 1, 1, 1, time.UTC)
	dateOfIncorporation := time.Date(2021, 9, 9, 1, 1, 1, 1, time.UTC)

	profile := entity2.Profile{
		ProfileID: &profileID,
		Person: entity2.Person{
			DocumentNumber: "12345678911222",
			PersonType:     values.PersonTypeCompany,
			PartnerID:      "Partner",
			Company: &entity2.Company{
				LegalName:            "Siderurgica",
				DateOfIncorporation:  &dateOfIncorporation,
				PlaceOfIncorporation: "BRA",
				ShareCapital: &entity2.MonetaryAmount{
					Amount: 20000,
				},
			},
		},
		CreatedAt: createdDate,
	}

	accounts := []entity2.Account{{AccountNumber: "xxx"}}
	addresses := []entity2.Address{}

	translatedAccounts := message.Accounts{{AgencyNumber: "123"}, {AgencyNumber: "456"}}
	translatedAddresses := message.Addresses{{}}
	translatedLegalName := "Siderurgica Nacional"
	translatedDateOfIncorporation := "01012010"
	translatedPlaceOfIncorporation := "BRA"
	translatedIncome := 20000.0
	date := time.Now().UTC()

	timeTranslatorInstance.On("TranslateTime", profile.CreatedAt).Return("01012020")
	timeTranslatorInstance.On("GenerateMessageDate").Return(date)
	accountTranslatorInstance.On("Translate", accounts).Return(translatedAccounts)
	addressTranslatorInstance.On("Translate", addresses).Return(translatedAddresses)
	companyTranslatorInstance.On("TranslateLegalName", profile.Company).Return(translatedLegalName)
	companyTranslatorInstance.On("TranslateDateOfIncorporation", profile.Company).Return(translatedDateOfIncorporation)
	companyTranslatorInstance.On("TranslatePlaceOfIncorporation", profile.Company).Return(translatedPlaceOfIncorporation)
	companyTranslatorInstance.On("TranslateIncome", profile.Company).Return(translatedIncome)

	expected := &message.TreeAdapterMessage{
		ProfileID:          profile.ProfileID.String(),
		Partner:            profile.PartnerID,
		ProfileCreatedDate: "01012020",
		Negotiation:        types.NegotiationTypeMesa,
		Person: message.Person{
			ID:                     profile.DocumentNumber,
			Type:                   types.PersonTypeLegalEntity,
			Name:                   translatedLegalName,
			BirthDate:              translatedDateOfIncorporation,
			Nationality:            translatedPlaceOfIncorporation,
			DigitalSign:            false,
			Addresses:              translatedAddresses,
			Accounts:               translatedAccounts,
			Contacts:               message.Contacts{},
			NotificationRecipients: message.NotificationRecipients{},
			LegalEntity: message.LegalEntity{
				Size:             types.CompanySizeNotInformed,
				BusinessActivity: types.BusinessActivityA,
				MonthlyIncome:    translatedIncome,
			},
		},
		Date: date,
	}

	received := translator.Translate(profile, accounts, addresses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

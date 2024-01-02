package limitMessageTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit/message"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := New()

	profileID := uuid.New()
	partnerID := uuid.New().String()

	profile := entity2.Profile{
		ProfileID: &profileID,
		Person: entity2.Person{
			DocumentNumber: "1234",
			PersonType:     values2.PersonTypeIndividual,
			PartnerID:      partnerID,
			OfferType:      "TestOffer",
			RoleType:       values2.RoleTypeCustomer,
		},
	}

	approvedRuleName := values2.RuleNameNotFoundInSerasa
	notApprovedRuleName := values2.RuleNameHasProblemsInSerasa

	state := entity2.State{
		ValidationStepsResults: []entity2.ValidationStepResult{
			{
				RuleResults: []entity2.RuleResultV2{
					{
						Result:   values2.ResultStatusApproved,
						RuleName: approvedRuleName,
					},
					{
						Result:   values2.ResultStatusAnalysing,
						RuleName: notApprovedRuleName,
					},
				},
			},
		},
	}

	expected := &message.LimitMessage{
		EventType:      message.EventTypeProfileApproved,
		ProfileID:      profile.ProfileID.String(),
		DocumentNumber: profile.DocumentNumber,
		PartnerID:      profile.PartnerID,
		OfferType:      profile.OfferType,
		PersonType:     string(profile.Person.PersonType),
		RoleType:       string(profile.RoleType),
		ApprovedRules:  []string{approvedRuleName.ToString()},
		Documents:      []string{},
	}

	received := translator.Translate(profile, state)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslate_with_Documents(t *testing.T) {
	translator := New()

	profileID := uuid.New()
	partnerID := uuid.New().String()

	profile := entity2.Profile{
		ProfileID: &profileID,
		Person: entity2.Person{
			DocumentNumber: "1234",
			PersonType:     values2.PersonTypeIndividual,
			PartnerID:      partnerID,
			OfferType:      "TestOffer",
			RoleType:       values2.RoleTypeCustomer,
			Documents: []entity2.Document{
				{
					DocumentType: values2.DocumentTypeIdentification,
				},
				{
					DocumentType: values2.DocumentTypeConstitutionDocument,
				},
			},
		},
	}

	approvedRuleName := values2.RuleNameDocumentNotFound
	notApprovedRuleName := values2.RuleNameHasProblemsInSerasa

	state := entity2.State{
		ValidationStepsResults: []entity2.ValidationStepResult{
			{
				RuleResults: []entity2.RuleResultV2{
					{
						Result:            values2.ResultStatusApproved,
						RuleName:          approvedRuleName,
						ApprovedDocuments: []string{"DOC_1"},
					},
					{
						Result:   values2.ResultStatusAnalysing,
						RuleName: notApprovedRuleName,
					},
				},
			},
		},
	}

	expected := &message.LimitMessage{
		EventType:      message.EventTypeProfileApproved,
		ProfileID:      profile.ProfileID.String(),
		DocumentNumber: profile.DocumentNumber,
		PartnerID:      profile.PartnerID,
		OfferType:      profile.OfferType,
		PersonType:     string(profile.Person.PersonType),
		RoleType:       string(profile.RoleType),
		ApprovedRules:  []string{approvedRuleName.ToString()},
		Documents: []string{profile.Documents[0].DocumentType,
			profile.Documents[1].DocumentType},
	}

	received := translator.Translate(profile, state)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

package profile

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAnalyze_should_approve_LR_rule(t *testing.T) {
	profileID := uuid.New()
	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: "partnerIdXxx",
			ProfileID: profileID,
			OfferType: "offerTypeXxx",
		}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	legalRepresentatives := []entity.LegalRepresentative{
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
	}

	profile.LegalRepresentatives = legalRepresentatives

	legalRepresentativeRule := NewLegalRepresentativeRule(profile, personProcessor)

	state0 := &entity.State{Result: values.ResultStatusApproved}

	personProcessor.On("ExecuteForPerson", legalRepresentatives[0].Person, profile.Person.OfferType).Return(state0, nil)

	expected := []entity.RuleResultV2{{
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetLegalRepresentatives,
		RuleName: values.RuleNameLegalRepresentativesResult,
		Pending:  false,
	},
	}

	received, err := legalRepresentativeRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

func TestAnalyze_should_mark_as_analysing_LR_rule(t *testing.T) {
	profileID := uuid.New()
	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: "partnerIdXxx",
			ProfileID: profileID,
			OfferType: "offerTypeXxx",
		}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	legalRepresentatives := []entity.LegalRepresentative{
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
	}

	profile.LegalRepresentatives = legalRepresentatives

	legalRepresentativeRule := NewLegalRepresentativeRule(profile, personProcessor)

	state0 := &entity.State{Result: values.ResultStatusAnalysing}

	personProcessor.On("ExecuteForPerson", legalRepresentatives[0].Person, profile.Person.OfferType).Return(state0, nil)

	expected := []entity.RuleResultV2{{
		Result:   values.ResultStatusAnalysing,
		RuleSet:  values.RuleSetLegalRepresentatives,
		RuleName: values.RuleNameLegalRepresentativesResult,
		Pending:  true,
		Metadata: []byte("[\"Legal Representative 00000000-0000-0000-0000-000000000000 is not Approved\"]"),
		Problems: []entity.Problem{{
			Code:   values.ProblemCodeLegalRepresentativeNotApproved,
			Detail: []string{"00000000-0000-0000-0000-000000000000"},
		}},
	},
	}

	received, err := legalRepresentativeRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

func TestAnalyze_should_reject_LR_rule(t *testing.T) {
	profileID := uuid.New()
	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: "partnerIdXxx",
			ProfileID: profileID,
			OfferType: "offerTypeXxx",
		}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	legalRepresentatives := []entity.LegalRepresentative{
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
	}

	profile.LegalRepresentatives = legalRepresentatives

	legalRepresentativeRule := NewLegalRepresentativeRule(profile, personProcessor)

	state0 := &entity.State{Result: values.ResultStatusRejected}

	personProcessor.On("ExecuteForPerson", legalRepresentatives[0].Person, profile.Person.OfferType).Return(state0, nil)

	expected := []entity.RuleResultV2{{
		Result:   values.ResultStatusRejected,
		RuleSet:  values.RuleSetLegalRepresentatives,
		RuleName: values.RuleNameLegalRepresentativesResult,
		Pending:  false,
		Metadata: []byte("[\"Legal Representative 00000000-0000-0000-0000-000000000000 is not Approved\"]"),
		Problems: []entity.Problem{{
			Code:   values.ProblemCodeLegalRepresentativeNotApproved,
			Detail: []string{"00000000-0000-0000-0000-000000000000"},
		}},
	},
	}

	received, err := legalRepresentativeRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

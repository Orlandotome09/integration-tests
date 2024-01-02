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

func TestAnalyze_should_approve_board_of_directors_rule(t *testing.T) {
	profileID := uuid.New()
	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: "partnerIdXxx",
			ProfileID: profileID,
			OfferType: "offerTypeXxx",
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedCompany: entity.EnrichedCompany{
					LegalNature: "2143",
				},
			},
		}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	boardOfDirectors := []entity.Director{
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
	}

	profile.BoardOfDirectors = boardOfDirectors

	boardOfDirectorsRule := NewBoardOfDirectorsRule(profile, personProcessor)

	state0 := &entity.State{Result: values.ResultStatusApproved}
	state1 := &entity.State{Result: values.ResultStatusApproved}

	personProcessor.On("ExecuteForPerson", boardOfDirectors[0].Person, profile.Person.OfferType).Return(state0, nil)
	personProcessor.On("ExecuteForPerson", boardOfDirectors[1].Person, profile.Person.OfferType).Return(state1, nil)

	expected := []entity.RuleResultV2{{
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetBoardOfDirectors,
		RuleName: values.RuleNameBoardOfDirectorsComplete,
		Pending:  false,
	}, {
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetBoardOfDirectors,
		RuleName: values.RuleNameBoardOfDirectorsResult,
		Pending:  false,
	},
	}

	received, err := boardOfDirectorsRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

func TestAnalyze_should_mark_as_analysing_board_of_directors_rule(t *testing.T) {
	profileID := uuid.New()
	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: "partnerIdXxx",
			ProfileID: profileID,
			OfferType: "offerTypeXxx",
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedCompany: entity.EnrichedCompany{
					LegalNature: "2143",
				},
			},
		}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	boardOfDirectors := []entity.Director{
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
	}

	profile.BoardOfDirectors = boardOfDirectors

	boardOfDirectorsRule := NewBoardOfDirectorsRule(profile, personProcessor)

	state0 := &entity.State{Result: values.ResultStatusApproved}
	state1 := &entity.State{Result: values.ResultStatusAnalysing}

	personProcessor.On("ExecuteForPerson", boardOfDirectors[0].Person, profile.Person.OfferType).Return(state0, nil)
	personProcessor.On("ExecuteForPerson", boardOfDirectors[1].Person, profile.Person.OfferType).Return(state1, nil)

	expected := []entity.RuleResultV2{{
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetBoardOfDirectors,
		RuleName: values.RuleNameBoardOfDirectorsComplete,
		Pending:  false,
	}, {
		Result:   values.ResultStatusAnalysing,
		RuleSet:  values.RuleSetBoardOfDirectors,
		RuleName: values.RuleNameBoardOfDirectorsResult,
		Pending:  true,
		Metadata: []byte("[\"Director 00000000-0000-0000-0000-000000000000 is not Approved\"]"),
		Problems: []entity.Problem{{
			Code:   values.ProblemCodeDirectorNotApproved,
			Detail: []string{"00000000-0000-0000-0000-000000000000"},
		}},
	},
	}

	received, err := boardOfDirectorsRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

func TestAnalyze_should_reject_board_of_directors_rule(t *testing.T) {
	profileID := uuid.New()
	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: "partnerIdXxx",
			ProfileID: profileID,
			OfferType: "offerTypeXxx",
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedCompany: entity.EnrichedCompany{
					LegalNature: "2143",
				},
			},
		}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	boardOfDirectors := []entity.Director{
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
		{Person: entity.Person{
			EntityID:  uuid.New(),
			PartnerID: profile.Person.PartnerID,
			ProfileID: profile.Person.ProfileID,
			OfferType: profile.Person.OfferType,
		}},
	}

	profile.BoardOfDirectors = boardOfDirectors

	boardOfDirectorsRule := NewBoardOfDirectorsRule(profile, personProcessor)

	state0 := &entity.State{Result: values.ResultStatusApproved}
	state1 := &entity.State{Result: values.ResultStatusRejected}

	personProcessor.On("ExecuteForPerson", boardOfDirectors[0].Person, profile.Person.OfferType).Return(state0, nil)
	personProcessor.On("ExecuteForPerson", boardOfDirectors[1].Person, profile.Person.OfferType).Return(state1, nil)

	expected := []entity.RuleResultV2{{
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetBoardOfDirectors,
		RuleName: values.RuleNameBoardOfDirectorsComplete,
		Pending:  false,
	}, {
		Result:   values.ResultStatusRejected,
		RuleSet:  values.RuleSetBoardOfDirectors,
		RuleName: values.RuleNameBoardOfDirectorsResult,
		Pending:  false,
		Metadata: []byte("[\"Director 00000000-0000-0000-0000-000000000000 is not Approved\"]"),
		Problems: []entity.Problem{{
			Code:   values.ProblemCodeDirectorNotApproved,
			Detail: []string{"00000000-0000-0000-0000-000000000000"},
		}},
	},
	}

	received, err := boardOfDirectorsRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

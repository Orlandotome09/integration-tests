package shareholders

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAnalyze_shouldFindShareholdersOK(t *testing.T) {
	profile := entity.Profile{Person: entity.Person{
		PartnerID: "partnerIdXxx",
		ProfileID: uuid.New(),
		OfferType: "offerTypeXxx",
	}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	shareholdersRule := NewShareholdersRule(profile, personProcessor)

	ownershipStructure := entity.OwnershipStructure{
		Shareholders: []entity.Shareholder{
			{Person: entity.Person{EntityID: uuid.New()}},
			{Person: entity.Person{EntityID: uuid.New()}},
		},
	}

	shareholder0 := ownershipStructure.Shareholders[0]
	shareholder0.Person.PartnerID = profile.Person.PartnerID
	shareholder0.Person.ProfileID = profile.Person.ProfileID
	shareholder0.Person.OfferType = profile.Person.OfferType

	shareholder1 := ownershipStructure.Shareholders[1]
	shareholder1.Person.PartnerID = profile.Person.PartnerID
	shareholder1.Person.ProfileID = profile.Person.ProfileID
	shareholder1.Person.OfferType = profile.Person.OfferType

	state0 := &entity.State{Result: values.ResultStatusApproved}
	state1 := &entity.State{Result: values.ResultStatusApproved}

	personProcessor.On("ExecuteForPerson", shareholder0.Person, profile.Person.OfferType).Return(state0, nil)
	personProcessor.On("ExecuteForPerson", shareholder1.Person, profile.Person.OfferType).Return(state1, nil)

	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholders,
		Pending:  false,
	}

	received, err := shareholdersRule.Analyze(ownershipStructure)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

func TestAnalyze_shouldFindShareholdersAnalysing(t *testing.T) {
	profile := entity.Profile{Person: entity.Person{
		PartnerID: "partnerIdXxx",
		ProfileID: uuid.New(),
		OfferType: "offerTypeXxx",
	}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	shareholdersRule := NewShareholdersRule(profile, personProcessor)

	shareholderNotApprovedID := uuid.New()
	shareholderNotApproved := entity.Shareholder{
		ShareholderID: &shareholderNotApprovedID,
		Person:        entity.Person{EntityID: uuid.New(), DocumentNumber: "xxx"},
	}
	ownershipStructure := entity.OwnershipStructure{
		Shareholders: []entity.Shareholder{
			{Person: entity.Person{EntityID: uuid.New()}},
			shareholderNotApproved,
		},
	}

	shareholder0 := ownershipStructure.Shareholders[0]
	shareholder0.Person.PartnerID = profile.Person.PartnerID
	shareholder0.Person.ProfileID = profile.Person.ProfileID
	shareholder0.Person.OfferType = profile.Person.OfferType

	shareholder1 := ownershipStructure.Shareholders[1]
	shareholder1.Person.PartnerID = profile.Person.PartnerID
	shareholder1.Person.ProfileID = profile.Person.ProfileID
	shareholder1.Person.OfferType = profile.Person.OfferType

	state0 := &entity.State{Result: values.ResultStatusApproved}
	stateNotApproved := &entity.State{Result: values.ResultStatusAnalysing}

	personProcessor.On("ExecuteForPerson", shareholder0.Person, profile.Person.OfferType).Return(state0, nil)
	personProcessor.On("ExecuteForPerson", shareholder1.Person, profile.Person.OfferType).Return(stateNotApproved, nil)

	notApprovedShareholders := []map[string]interface{}{
		{
			"document_number": shareholderNotApproved.Person.DocumentNumber,
			"shareholder_id":  shareholderNotApproved.ShareholderID.String(),
		},
	}
	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusAnalysing,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholders,
		Pending:  true,
		Problems: []entity.Problem{
			{
				Code:   values.ProblemCodeShareholderNotApproved,
				Detail: notApprovedShareholders,
			},
		},
	}

	received, err := shareholdersRule.Analyze(ownershipStructure)

	assert.Nil(t, err)
	assert.Equal(t, expected.Result, received.Result)
	assert.Equal(t, expected.Pending, received.Pending)
	assert.Equal(t, expected.Problems, received.Problems)
	assert.NotNil(t, received.Metadata)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

func TestAnalyze_shouldFindShareholdersRejected(t *testing.T) {
	profile := entity.Profile{Person: entity.Person{
		PartnerID: "partnerIdXxx",
		ProfileID: uuid.New(),
		OfferType: "offerTypeXxx",
	}}
	personProcessor := &mocks.CompliancePersonProcessor{}

	shareholdersRule := NewShareholdersRule(profile, personProcessor)

	shareholderNotApprovedID := uuid.New()
	shareholderNotApproved := entity.Shareholder{
		ShareholderID: &shareholderNotApprovedID,
		Person:        entity.Person{EntityID: uuid.New(), DocumentNumber: "xxx"},
	}
	ownershipStructure := entity.OwnershipStructure{
		Shareholders: []entity.Shareholder{
			{Person: entity.Person{EntityID: uuid.New()}},
			shareholderNotApproved,
		},
	}

	shareholder0 := ownershipStructure.Shareholders[0]
	shareholder0.Person.PartnerID = profile.Person.PartnerID
	shareholder0.Person.ProfileID = profile.Person.ProfileID
	shareholder0.Person.OfferType = profile.Person.OfferType

	shareholder1 := ownershipStructure.Shareholders[1]
	shareholder1.Person.PartnerID = profile.Person.PartnerID
	shareholder1.Person.ProfileID = profile.Person.ProfileID
	shareholder1.Person.OfferType = profile.Person.OfferType

	state0 := &entity.State{Result: values.ResultStatusApproved}
	stateNotApproved := &entity.State{Result: values.ResultStatusRejected}

	personProcessor.On("ExecuteForPerson", shareholder0.Person, profile.Person.OfferType).Return(state0, nil)
	personProcessor.On("ExecuteForPerson", shareholder1.Person, profile.Person.OfferType).Return(stateNotApproved, nil)

	notApprovedShareholders := []map[string]interface{}{
		{
			"document_number": shareholderNotApproved.Person.DocumentNumber,
			"shareholder_id":  shareholderNotApproved.ShareholderID.String(),
		},
	}
	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusRejected,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholders,
		Pending:  false,
		Problems: []entity.Problem{
			{
				Code:   values.ProblemCodeShareholderNotApproved,
				Detail: notApprovedShareholders,
			},
		},
	}

	received, err := shareholdersRule.Analyze(ownershipStructure)

	assert.Nil(t, err)
	assert.Equal(t, expected.Result, received.Result)
	assert.Equal(t, expected.Pending, received.Pending)
	assert.Equal(t, expected.Problems, received.Problems)
	assert.NotNil(t, received.Metadata)
	mock.AssertExpectationsForObjects(t, personProcessor)
}

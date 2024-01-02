package shareholding

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyze_shouldFindMinimumShareholdingInEnrichedOwnershipStructure(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			DocumentNumber: "12456",
			OfferType:      "xxx",
			PartnerID:      "partnerId",
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						RulesConfig: &entity.RuleSetConfig{
							OwnershipStructureParams: &entity.OwnershipStructureParams{},
						},
					},
				},
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedCompany: entity.EnrichedCompany{
					OwnershipStructure: &entity.OwnershipStructure{ShareholdingSum: 96.0},
				},
			},
		},
	}
	shareholdingRule := NewShareholdingRule(profile)

	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholding,
		Pending:  false,
	}
	expectedOwnershipStructure := profile.EnrichedInformation.OwnershipStructure

	received, receivedOwnershipStructure, err := shareholdingRule.Analyze()

	assert.Equal(t, expected.Result, received.Result)
	assert.Equal(t, expectedOwnershipStructure, receivedOwnershipStructure)
	assert.Nil(t, err)
}

func TestAnalyze_shouldFindMinimumShareholdingInManuallyFilledOwnershipStructure(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			DocumentNumber: "12456",
			OfferType:      "xxx",
			PartnerID:      "partnerId",
			EntityID:       uuid.New(),
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						RulesConfig: &entity.RuleSetConfig{
							OwnershipStructureParams: &entity.OwnershipStructureParams{},
						},
					},
				},
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedCompany: entity.EnrichedCompany{
					OwnershipStructure: &entity.OwnershipStructure{ShareholdingSum: 94.0},
				},
			},
		},
		OwnershipStructure: &entity.OwnershipStructure{ShareholdingSum: 96.0},
	}

	shareholdingRule := NewShareholdingRule(profile)

	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusApproved,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholding,
		Pending:  false,
	}
	expectedOwnershipStructure := profile.OwnershipStructure

	received, receivedOwnershipStructure, err := shareholdingRule.Analyze()

	assert.Equal(t, expected.Result, received.Result)
	assert.Equal(t, expectedOwnershipStructure, receivedOwnershipStructure)
	assert.Nil(t, err)
}

func TestAnalyze_shouldNotFindMinimumShareholdingNeitherInEnrichedNorInFilledOwnershipStructure(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			DocumentNumber: "12456",
			OfferType:      "xxx",
			PartnerID:      "partnerId",
			EntityID:       uuid.New(),
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						RulesConfig: &entity.RuleSetConfig{
							OwnershipStructureParams: &entity.OwnershipStructureParams{},
						},
					},
				},
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedCompany: entity.EnrichedCompany{
					OwnershipStructure: &entity.OwnershipStructure{ShareholdingSum: 94.0},
				},
			},
		},
		OwnershipStructure: &entity.OwnershipStructure{ShareholdingSum: 90.0},
	}

	shareholdingRule := NewShareholdingRule(profile)

	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusAnalysing,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholding,
		Pending:  true,
		Problems: []entity.Problem{
			{
				Code: values.ProblemCodeShareholdingNotAchieveMinimumRequired,
				Detail: map[string]interface{}{
					"min_shareholding":    95.0,
					"enriched_percentage": profile.EnrichedInformation.OwnershipStructure.ShareholdingSum,
					"manually_percentage": profile.OwnershipStructure.ShareholdingSum,
				},
			},
		},
	}
	expectedOwnershipStructure := profile.OwnershipStructure

	received, receivedOwnershipStructure, err := shareholdingRule.Analyze()

	assert.Equal(t, expected.Result, received.Result)
	assert.Equal(t, expected.Pending, received.Pending)
	assert.Equal(t, expected.Problems, received.Problems)
	assert.NotNil(t, received.Metadata)
	assert.Equal(t, expectedOwnershipStructure, receivedOwnershipStructure)
	assert.Nil(t, err)
}

func TestAnalyze_shouldNotFindFilledOwnershipStructure(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			DocumentNumber: "12456",
			OfferType:      "xxx",
			PartnerID:      "partnerId",
			EntityID:       uuid.New(),
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						RulesConfig: &entity.RuleSetConfig{
							OwnershipStructureParams: &entity.OwnershipStructureParams{},
						},
					},
				},
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedCompany: entity.EnrichedCompany{
					OwnershipStructure: &entity.OwnershipStructure{ShareholdingSum: 94.0},
				},
			},
		},
	}

	shareholdingRule := NewShareholdingRule(profile)

	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusAnalysing,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholding,
		Pending:  true,
		Problems: []entity.Problem{
			{
				Code: values.ProblemCodeShareholdingNotAchieveMinimumRequired,
				Detail: map[string]interface{}{
					"min_shareholding":    95.0,
					"enriched_percentage": profile.EnrichedInformation.OwnershipStructure.ShareholdingSum,
					"manually_percentage": 0.0,
				},
			},
		},
	}
	expectedOwnershipStructure := profile.EnrichedInformation.OwnershipStructure

	received, receivedOwnershipStructure, err := shareholdingRule.Analyze()

	assert.Equal(t, expected.Result, received.Result)
	assert.Equal(t, expected.Pending, received.Pending)
	assert.Equal(t, expected.Problems, received.Problems)
	assert.NotNil(t, received.Metadata)
	assert.Equal(t, expectedOwnershipStructure, receivedOwnershipStructure)
	assert.Nil(t, err)
}

func TestAnalyze_shouldNotFindEnrichedOwnershipStructureNeitherFilledOwnershipStructure(t *testing.T) {
	profile := entity.Profile{Person: entity.Person{
		DocumentNumber: "12456",
		OfferType:      "xxx",
		PartnerID:      "partnerId",
		EntityID:       uuid.New(),
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						OwnershipStructureParams: &entity.OwnershipStructureParams{},
					},
				},
			},
		},
	}}

	shareholdingRule := NewShareholdingRule(profile)

	expected := &entity.RuleResultV2{
		Result:   values.ResultStatusAnalysing,
		RuleSet:  values.RuleSetOwnershipStructure,
		RuleName: values.RuleNameShareholding,
		Pending:  true,
		Problems: []entity.Problem{
			{
				Code: values.ProblemCodeShareholdingNotAchieveMinimumRequired,
				Detail: map[string]interface{}{
					"min_shareholding":    95.0,
					"enriched_percentage": 0.0,
					"manually_percentage": 0.0,
				},
			},
		},
	}
	var expectedOwnershipStructure *entity.OwnershipStructure = nil

	received, receivedOwnershipStructure, err := shareholdingRule.Analyze()

	assert.Equal(t, expected.Result, received.Result)
	assert.Equal(t, expected.Pending, received.Pending)
	assert.Equal(t, expected.Problems, received.Problems)
	assert.NotNil(t, received.Metadata)
	assert.Equal(t, expectedOwnershipStructure, receivedOwnershipStructure)
	assert.Nil(t, err)
}

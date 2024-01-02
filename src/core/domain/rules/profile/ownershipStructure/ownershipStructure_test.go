package ownershipStructure

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAnalyze_shouldValidateShareholdersByApprovingMinimumShareholding(t *testing.T) {
	profile := entity.Profile{}
	shareholdingRule := &mocks.ShareholdingAnalyzer{}
	shareholdersRule := &mocks.ShareholdersAnalyzer{}
	ownershipStructureRule := NewOwnershipStructureRule(profile, shareholdingRule, shareholdersRule)

	shareholdingResult := &entity.RuleResultV2{Result: values.ResultStatusApproved}
	ownershipStructure := &entity.OwnershipStructure{ShareholdingSum: 99.0}
	shareholdersResult := &entity.RuleResultV2{Result: values.ResultStatusApproved}

	shareholdingRule.On("Analyze").Return(shareholdingResult, ownershipStructure, nil)
	shareholdersRule.On("Analyze", *ownershipStructure).Return(shareholdersResult, nil)

	expected := []entity.RuleResultV2{*shareholdingResult, *shareholdersResult}
	received, err := ownershipStructureRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, shareholdingRule, shareholdersRule)
}

func TestAnalyze_shouldValidateShareholdersByFindingOverride(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			Overrides: []entity.Override{
				{
					Result:   values.ResultStatusApproved,
					RuleSet:  values.RuleSetOwnershipStructure,
					RuleName: values.RuleNameShareholding,
				},
			},
		},
	}
	shareholdingRule := &mocks.ShareholdingAnalyzer{}
	shareholdersRule := &mocks.ShareholdersAnalyzer{}
	ownershipStructureRule := NewOwnershipStructureRule(profile, shareholdingRule, shareholdersRule)

	shareholdingResult := &entity.RuleResultV2{Result: values.ResultStatusAnalysing}

	ownershipStructure := &entity.OwnershipStructure{ShareholdingSum: 99.0}
	shareholdersResult := &entity.RuleResultV2{Result: values.ResultStatusApproved}

	shareholdingRule.On("Analyze").Return(shareholdingResult, ownershipStructure, nil)
	shareholdersRule.On("Analyze", *ownershipStructure).Return(shareholdersResult, nil)

	expected := []entity.RuleResultV2{*shareholdingResult, *shareholdersResult}
	received, err := ownershipStructureRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, shareholdingRule, shareholdersRule)
}

func TestAnalyze_shouldNotValidateShareholders(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			Overrides: []entity.Override{
				{
					Result:   values.ResultStatusAnalysing,
					RuleSet:  values.RuleSetOwnershipStructure,
					RuleName: values.RuleNameShareholding,
				},
			},
		},
	}
	shareholdingRule := &mocks.ShareholdingAnalyzer{}
	shareholdersRule := &mocks.ShareholdersAnalyzer{}
	ownershipStructureRule := NewOwnershipStructureRule(profile, shareholdingRule, shareholdersRule)

	shareholdingResult := &entity.RuleResultV2{Result: values.ResultStatusAnalysing}

	ownershipStructure := &entity.OwnershipStructure{ShareholdingSum: 99.0}
	shareholdersResult := entity.NewRuleResultV2(values.RuleSetOwnershipStructure, values.RuleNameShareholders).SetResult(values.ResultStatusIgnored)

	shareholdingRule.On("Analyze").Return(shareholdingResult, ownershipStructure, nil)

	expected := []entity.RuleResultV2{*shareholdingResult, *shareholdersResult}
	received, err := ownershipStructureRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
	mock.AssertExpectationsForObjects(t, shareholdingRule, shareholdersRule)
	shareholdersRule.AssertNumberOfCalls(t, "Analyze", 0)
	mock.AssertExpectationsForObjects(t, shareholdingRule)
}

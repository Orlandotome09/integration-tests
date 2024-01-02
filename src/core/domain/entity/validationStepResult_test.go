package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindStepApprovedRules(t *testing.T) {
	step := ValidationStepResult{
		RuleResults: []RuleResultV2{
			{RuleName: values.RuleNameBlocked, Result: values.ResultStatusApproved},
			{RuleName: values.RuleNameShareholders, Result: values.ResultStatusAnalysing},
			{RuleName: values.RuleNameShareholding, Result: values.ResultStatusApproved},
			{RuleName: values.RuleNameAddressNotFound, Result: values.ResultStatusAnalysing},
		},
	}

	expected := []values.RuleName{
		values.RuleNameBlocked,
		values.RuleNameShareholding,
	}

	received := step.FindApprovedRules()

	assert.Equal(t, expected, received)
}

func Test_FindStepsApprovedRules(t *testing.T) {
	steps := ValidationStepsResult{
		ValidationStepResult{
			RuleResults: []RuleResultV2{
				{RuleName: values.RuleNameBlocked, Result: values.ResultStatusApproved},
				{RuleName: values.RuleNameShareholders, Result: values.ResultStatusAnalysing},
				{RuleName: values.RuleNameShareholding, Result: values.ResultStatusApproved},
				{RuleName: values.RuleNameAddressNotFound, Result: values.ResultStatusAnalysing},
			},
		},
		ValidationStepResult{
			RuleResults: []RuleResultV2{
				{RuleName: values.RuleNameBlocked, Result: values.ResultStatusAnalysing},
				{RuleName: values.RuleNameShareholders, Result: values.ResultStatusApproved},
				{RuleName: values.RuleNameShareholding, Result: values.ResultStatusAnalysing},
				{RuleName: values.RuleNameAddressNotFound, Result: values.ResultStatusApproved},
			},
		},
	}

	expected := []values.RuleName{
		values.RuleNameBlocked,
		values.RuleNameShareholding,
		values.RuleNameShareholders,
		values.RuleNameAddressNotFound,
	}

	received := steps.FindApprovedRules()

	assert.Equal(t, expected, received)
}

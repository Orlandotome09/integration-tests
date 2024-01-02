package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type ValidationStepResult struct {
	StepNumber      int            `json:"step_number"`
	Result          values.Result  `json:"result"`
	SkipForApproval bool           `json:"skip_for_approval"`
	RuleResults     []RuleResultV2 `json:"rule_results"`
}

func (validationStepResult ValidationStepResult) FindApprovedRules() []values.RuleName {
	approvedRules := make([]values.RuleName, 0)
	for _, rule := range validationStepResult.RuleResults {
		if rule.IsApproved() {
			approvedRules = append(approvedRules, rule.RuleName)
		}
	}

	return approvedRules
}

type ValidationStepsResult []ValidationStepResult

func (validationStepsResult ValidationStepsResult) FindApprovedRules() []values.RuleName {
	approvedRules := make([]values.RuleName, 0)
	for _, validationStep := range validationStepsResult {
		approvedRules = append(approvedRules, validationStep.FindApprovedRules()...)
	}

	return approvedRules
}

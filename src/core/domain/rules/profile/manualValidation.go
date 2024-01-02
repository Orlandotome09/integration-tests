package profile

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type manualValidationRule struct {
	profile entity.Profile
}

func NewManualValidationRule(profile entity.Profile) entity.Rule {
	return &manualValidationRule{
		profile: profile,
	}
}

func (ref *manualValidationRule) Analyze() ([]entity.RuleResultV2, error) {
	manualValidationRule := entity.NewRuleResultV2(values.RuleSetManualValidation, values.RuleNameManualValidation)
	manualValidationRule.SetResult(values.ResultStatusAnalysing).SetPending(true)

	return []entity.RuleResultV2{*manualValidationRule}, nil
}

func (ref *manualValidationRule) Name() values.RuleSet {
	return values.RuleSetManualValidation
}

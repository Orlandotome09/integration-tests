package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type manualBlockAnalyzer struct {
	person entity.Person
}

func NewManualBlockAnalyzer(person entity.Person) entity.Rule {
	return &manualBlockAnalyzer{
		person: person,
	}
}

func (ref *manualBlockAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	manualBlockRule := entity.NewRuleResultV2(values.RuleSetManualBlock, values.RuleNameManualBlock)
	manualBlockRule.SetResult(values.ResultStatusApproved)

	return []entity.RuleResultV2{*manualBlockRule}, nil
}

func (ref *manualBlockAnalyzer) Name() values.RuleSet {
	return values.RuleSetManualBlock
}

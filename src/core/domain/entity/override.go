package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"time"

	"github.com/google/uuid"
)

type Override struct {
	EntityID   uuid.UUID         `json:"entity_id"`
	EntityType values.EntityType `json:"entity_type"`
	ParentID   *uuid.UUID        `json:"parent_id"`
	RuleSet    values.RuleSet    `json:"rule_set"`
	RuleName   values.RuleName   `json:"rule_name"`
	Result     values.Result     `json:"result"`
	Author     string            `json:"author"`
	Comments   string            `json:"comments"`
	CreatedAt  time.Time         `json:"created_at"`
}

func (override Override) Validate() error {
	err := override.EntityType.Validate()
	if err != nil {
		return err
	}

	err = override.RuleSet.Validate()
	if err != nil {
		return err
	}

	err = override.RuleName.Validate()
	if err != nil {
		return err
	}

	err = override.Result.Validate()
	if err != nil {
		return err
	}

	return nil
}

type Overrides []Override

func (overrides Overrides) FindByRuleSetAndName(ruleSet values.RuleSet, ruleName values.RuleName) (Override, bool) {
	for _, override := range overrides {
		if override.RuleSet == ruleSet && override.RuleName == ruleName {
			return override, true
		}
	}

	return Override{}, false
}

func (overrides Overrides) HasBlocked() bool {
	override, exists := overrides.FindByRuleSetAndName(values.RuleSetManualBlock, values.RuleNameManualBlock)
	return exists && override.Result == values.ResultStatusBlocked
}

func (overrides Overrides) HasInactive() bool {
	override, exists := overrides.FindByRuleSetAndName(values.RuleSetState, values.RuleNameInactive)
	return exists && override.Result == values.ResultStatusInactive
}

func (overrides Overrides) HasAnalizing() bool {
	override, exists := overrides.FindByRuleSetAndName(values.RuleSetState, values.RuleNameManualValidation)
	return exists && override.Result == values.ResultStatusAnalysing
}

package _interfaces

import (
	entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ProfileRulesFactory interface {
	GetRules(ruleSetConfig *entity.RuleSetConfig, profile *entity.Profile) []entity.Rule
}

type PersonRulesFactory interface {
	GetRules(ruleSetConfig entity.RuleSetConfig, person entity.Person) []entity.Rule
}

type ContractRulesFactory interface {
	GetRules(contract entity.Contract) []entity.Rule
}

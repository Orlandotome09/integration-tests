package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type Rule interface {
	Analyze() ([]RuleResultV2, error)
	Name() values.RuleSet
}

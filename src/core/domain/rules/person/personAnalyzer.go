package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type Analyzer interface {
	Analyze(person entity.Person) (*entity.RuleResultV2, error)
}

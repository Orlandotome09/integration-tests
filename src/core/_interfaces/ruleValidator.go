package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type RuleValidator interface {
	Validate(state entity.State, override entity.Overrides, noCache bool, entityID uuid.UUID, engineName string) (*entity.State, error)
	SetRules(rules []entity.RuleValidatorStep)
	NewInstance() RuleValidator
}

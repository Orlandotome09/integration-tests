package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type OverrideRepository interface {
	Save(override entity.Override) error
	Delete(entityID uuid.UUID, set values2.RuleSet, name values2.RuleName) error
	FindByEntityID(entityID uuid.UUID) (entity.Overrides, error)
}

type OverrideService interface {
	Save(override entity.Override) error
	Delete(override entity.Override) error
	FindByEntityID(entityID uuid.UUID) (entity.Overrides, error)
}

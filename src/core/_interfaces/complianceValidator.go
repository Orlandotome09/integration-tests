package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type ComplianceValidator interface {
	Validate(state entity.State, override entity.Overrides, noCache bool, entityID uuid.UUID, engineName string) (*entity.State, error)
	PosProcessing(previousState *entity.State, newState *entity.State, entityID uuid.UUID) error
	GetName() string
}

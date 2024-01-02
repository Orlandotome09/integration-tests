package _interfaces

import (
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type Engine interface {
	Prepare(entityID uuid.UUID) error
	NewInstance() Engine
	ComplianceValidator
}

type EngineSelector interface {
	GetEngineNameAndEntityId(event values2.Event) (values2.EngineName, *uuid.UUID, error)
}

type EngineFactory interface {
	CreateEngine(engineName string) (Engine, error)
}

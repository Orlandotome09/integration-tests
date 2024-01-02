package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type EventPublisher interface {
	SendEvent(entityID uuid.UUID, parentID *uuid.UUID, engineName values.EngineName, eventType values.EventType) (*entity.State, error)
}

type StateEventsPublisher interface {
	Send(state entity.State, eventType values.EventType) error
}

type ComplianceCommandPublisher interface {
	SendCommand(entityID uuid.UUID, parentID *uuid.UUID, engineName values.EngineName) (*entity.State, error)
}

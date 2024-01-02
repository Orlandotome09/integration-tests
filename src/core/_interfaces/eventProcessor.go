package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type EventProcessor interface {
	ExecuteForEvent(event *values.Event) (*entity.State, error)
}

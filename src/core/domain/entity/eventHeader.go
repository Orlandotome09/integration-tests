package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"time"
)

type EventHeader struct {
	ID         string            `json:"id"`
	EntityID   string            `json:"entity_id"`
	EntityType values.EntityType `json:"entity_type"`
	EventType  values.EventType  `json:"event_type"`
	UpdateDate time.Time         `json:"update_date"`
}

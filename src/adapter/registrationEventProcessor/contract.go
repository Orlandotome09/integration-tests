package registrationEventProcessor

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type RegistrationEvent struct {
	EventType  string          `json:"event_type"`
	ProfileID  uuid.UUID       `json:"profile_id"`
	EntityID   uuid.UUID       `json:"entity_id"`
	EntityType string          `json:"entity_type"`
	ParentID   *uuid.UUID      `json:"parent_id,omitempty"`
	ParentType string          `json:"parent_type"`
	UpdateDate time.Time       `json:"update_date"`
	Content    json.RawMessage `json:"content"`
}

func (event RegistrationEvent) IsValid() (bool, error) {
	if event.UpdateDate.IsZero() {
		return false, errors.New("invalid update date")
	}

	if event.EntityID == uuid.Nil {
		return false, errors.New("event without entity_id")
	}

	return true, nil
}

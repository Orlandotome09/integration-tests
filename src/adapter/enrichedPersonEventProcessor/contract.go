package enrichedPersonEventProcessor

import (
	"errors"

	"github.com/google/uuid"
)

type PersonEnrichedData struct {
	PersonID   uuid.UUID   `json:"person_id"`
	PersonType string      `json:"person_type"`
	Situation  int         `json:"situation"`
	Content    interface{} `json:"content"`
}

type EnrichedPersonEvent struct {
	EntityID   uuid.UUID          `json:"id"`
	EntityType string             `json:"entity_type"`
	EventType  string             `json:"event_type"`
	Data       PersonEnrichedData `json:"data"`
}

func (event EnrichedPersonEvent) IsValid() (bool, error) {
	if event.EntityID == uuid.Nil {
		return false, errors.New("event without profile id")
	}

	return true, nil
}

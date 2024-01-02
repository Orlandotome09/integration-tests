package complianceCommandProcessor

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type complianceCommand struct {
	EngineName string      `json:"engine_name"`
	EventType  string      `json:"event_type"`
	ParentID   *uuid.UUID  `json:"parent_id"`
	EntityID   uuid.UUID   `json:"entity_id"`
	Date       time.Time   `json:"date"`
	Content    interface{} `json:"content"`
}

func (command complianceCommand) IsValid() (bool, error) {
	if command.EntityID == uuid.Nil {
		return false, errors.New("event without profile id")
	}

	return true, nil
}

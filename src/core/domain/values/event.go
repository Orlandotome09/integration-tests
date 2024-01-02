package values

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	EngineName  EngineName
	EventType   EventType
	ParentID    uuid.UUID
	EntityID    uuid.UUID
	Date        time.Time
	RequestDate time.Time
	Content     interface{}
}

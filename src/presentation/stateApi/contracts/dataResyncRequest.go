package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type DataResyncRequest struct {
	Ids       []string         `json:"ids"`
	EventType values.EventType `json:"event_type"`
}

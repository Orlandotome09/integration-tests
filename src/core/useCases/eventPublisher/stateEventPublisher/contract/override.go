package contract

import (
	"time"

	"github.com/google/uuid"
)

type Override struct {
	EntityID   uuid.UUID  `json:"entity_id"`
	EntityType string     `json:"entity_type"`
	ParentID   *uuid.UUID `json:"parent_id"`
	RuleSet    string     `json:"rule_set"`
	RuleName   string     `json:"rule_name"`
	Result     string     `json:"result"`
	Author     string     `json:"author"`
	Comments   string     `json:"comments"`
	CreatedAt  time.Time  `json:"created_at"`
}

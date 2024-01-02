package entity

import "github.com/google/uuid"

type Shareholder struct {
	Person           Person     `json:"person"`
	ShareholderID    *uuid.UUID `json:"shareholder_id"`
	ParentID         *uuid.UUID `json:"parent_id"`
	Level            int        `json:"level"`
	OwnershipPercent float64    `json:"ownership_percent"`
}

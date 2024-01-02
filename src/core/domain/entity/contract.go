package entity

import (
	"github.com/google/uuid"
	"time"
)

type Contract struct {
	ContractID           *uuid.UUID `json:"contract_id"`
	EstimatedTotalAmount float64    `json:"estimated_total_amount"`
	DueTime              string     `json:"due_time"`
	Installments         int        `json:"installments"`
	CorrelationID        string     `json:"correlation_id"`
	ProfileID            *uuid.UUID `json:"profile_id"`
	DocumentID           *uuid.UUID `json:"document_id"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

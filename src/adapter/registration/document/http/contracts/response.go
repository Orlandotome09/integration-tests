package contracts

import (
	"github.com/google/uuid"
	"time"
)

type DocumentResponse struct {
	DocumentID     uuid.UUID      `json:"document_id"`
	EntityID       uuid.UUID      `json:"entity_id"`
	Type           string         `json:"type" binding:"required"`
	SubType        string         `json:"sub_type,omitempty"`
	DocumentFields DocumentFields `json:"document_fields,omitempty"`
	ExpirationDate string         `json:"expiration_date,omitempty"`
	EmissionDate   *time.Time     `json:"emission_date,omitempty"`
}

type DocumentFields struct {
	Number    string `json:"number,omitempty"`
	IssueDate string `json:"issue_date,omitempty"`
	Name      string `json:"name,omitempty"`
}

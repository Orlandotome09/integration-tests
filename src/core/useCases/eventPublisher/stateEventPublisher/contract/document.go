package contract

import (
	"github.com/google/uuid"
)

type Document struct {
	DocumentID      uuid.UUID      `json:"document_id"`
	EntityID        uuid.UUID      `json:"-"`
	DocumentType    string         `json:"document_type"`
	DocumentSubType string         `json:"document_sub_type"`
	DocumentFields  DocumentFields `json:"-"`
}

type DocumentFields struct {
	Number    string
	IssueDate string
	Name      string
}

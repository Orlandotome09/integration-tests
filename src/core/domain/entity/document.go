package entity

import (
	"github.com/google/uuid"
	"time"
)

type Document struct {
	DocumentID      uuid.UUID      `json:"document_id"`
	EntityID        uuid.UUID      `json:"-"`
	DocumentType    string         `json:"document_type"`
	DocumentSubType string         `json:"document_sub_type"`
	DocumentFields  DocumentFields `json:"-"`
	ExpirationDate  string         `json:"expiration_date,omitempty"`
	EmissionDate    *time.Time     `json:"emission_date,omitempty"`
}

type DocumentFields struct {
	Number    string
	IssueDate string
	Name      string
}

func (document Document) IsSameTypeOf(documentRequired DocumentRequired) bool {
	return document.DocumentType == documentRequired.DocumentType
}

func (document Document) IsSameSubtypeOf(documentRequired DocumentRequired) bool {
	return document.DocumentSubType == documentRequired.DocumentSubType
}

type Documents []Document

func (documents Documents) AtLeastOneHasFile(documentFiles DocumentFiles) bool {
	for _, document := range documents {
		if documentFiles.HaveDocument(document.DocumentID) {
			return true
		}

	}

	return false
}

func (documents Documents) HaveRequiredFile(documentFiles []DocumentFile, needFile bool) bool {
	return !needFile || documents.AtLeastOneHasFile(documentFiles)
}

func (documents Documents) HaveDocumentFile(documentFile DocumentFile) bool {
	for _, document := range documents {
		if document.DocumentID == documentFile.DocumentID {
			return true
		}
	}

	return false
}

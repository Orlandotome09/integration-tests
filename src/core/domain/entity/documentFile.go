package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"time"
)

type DocumentFile struct {
	DocumentFileID *uuid.UUID      `json:"document_file_id"`
	DocumentID     uuid.UUID       `json:"document_id"`
	FileID    uuid.UUID       `json:"file_id"`
	FileSide  values.FileSide `json:"file_side"`
	CreatedAt time.Time       `json:"created_at"`
}

type DocumentFiles []DocumentFile

func (documentFiles DocumentFiles) HaveDocument(documentID uuid.UUID) bool {
	for _, documentFile := range documentFiles {
		if documentFile.DocumentID == documentID {
			return true
		}
	}
	return false
}

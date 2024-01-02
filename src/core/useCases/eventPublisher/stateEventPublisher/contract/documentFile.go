package contract

import (
	"github.com/google/uuid"
	"time"
)

type DocumentFile struct {
	DocumentFileID *uuid.UUID `json:"document_file_id"`
	DocumentID     uuid.UUID  `json:"document_id"`
	FileID         uuid.UUID  `json:"file_id"`
	FileSide       string     `json:"file_side"`
	CreatedAt      time.Time  `json:"created_at"`
}

type DocumentFiles []DocumentFile

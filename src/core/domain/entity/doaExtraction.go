package entity

import (
	"github.com/google/uuid"
)

type DOAExtraction struct {
	Message   string    `json:"message"`
	RequestID uuid.UUID `json:"request_id"`
}

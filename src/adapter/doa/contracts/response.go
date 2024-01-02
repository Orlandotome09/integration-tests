package contracts

import (
	"github.com/google/uuid"
)

type DOAExtractionResponse struct {
	Message   string    `json:"message"`
	RequestID uuid.UUID `json:"request_id"`
}

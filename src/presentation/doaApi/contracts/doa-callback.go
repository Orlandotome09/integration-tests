package contracts

import (
	"github.com/google/uuid"
)

type DOACallback struct {
	RequestID uuid.UUID `json:"request_id"`
	Scores    Scores    `json:"scores"`
}

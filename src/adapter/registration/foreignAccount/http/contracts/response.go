package contracts

import (
	"github.com/google/uuid"
)

type ForeignAccountResponse struct {
	ForeignAccountID uuid.UUID `json:"foreign_account_id"`
	ProfileID        uuid.UUID `json:"profile_id"`
}

package entity

import (
	"github.com/google/uuid"
)

type ForeignAccount struct {
	ForeignAccountID uuid.UUID
	ProfileID        uuid.UUID
}

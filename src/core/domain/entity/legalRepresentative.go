package entity

import (
	"github.com/google/uuid"
	"time"
)

type LegalRepresentative struct {
	Person
	LegalRepresentativeID uuid.UUID  `json:"legal_representative_id"`
	ExpirationDate        *time.Time `json:"expiration_date,omitempty"`
}

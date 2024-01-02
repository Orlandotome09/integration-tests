package contracts

import (
	"github.com/google/uuid"
	"time"
)

type BoardOfDirectorsResponse struct {
	DirectorID     uuid.UUID `json:"director_id"`
	ProfileID      uuid.UUID `json:"profile_id"`
	FullName       string    `json:"full_name" binding:"required"`
	DocumentNumber string    `json:"document_number" binding:"required"`
	DateOfBirth    string    `json:"date_of_birth,omitempty"`
	Nationality    string    `json:"nationality"`
	Pep            *bool     `json:"pep,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

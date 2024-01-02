package contract

import "github.com/google/uuid"

type Contact struct {
	ContactID      *uuid.UUID `json:"contact_id"`
	ProfileID      *uuid.UUID `json:"profile_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	Phones         []Phone    `json:"phones"`
	Nationality    string     `json:"nationality"`
	DocumentNumber string     `json:"document_number"`
}

package contracts

import "github.com/google/uuid"

type ContactResponse struct {
	ContactID      uuid.UUID `json:"contact_id"`
	ProfileID      uuid.UUID `json:"profile_id"`
	Name           string    `json:"name"`
	DocumentNumber string    `json:"document_number"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Phones         Phones    `json:"phones,omitempty"`
	Nationality    string    `json:"nationality"`
}

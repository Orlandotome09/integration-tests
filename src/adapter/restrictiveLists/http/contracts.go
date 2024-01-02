package restrictiveListsHttpClient

import "time"

type InternalListResponse []InternalListItemResponse

type InternalListItemResponse struct {
	FullName       string    `json:"full_name"`
	DocumentNumber string    `json:"document_number"`
	PersonType     string    `json:"person_type"`
	Author         string    `json:"author"`
	Justification  string    `json:"justification"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PepResponse struct {
	DocumentNumber string    `json:"document_number"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	Institution    string    `json:"institution"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	Source         string    `json:"source"`
	CreatedAt      time.Time `json:"created_at"`
}

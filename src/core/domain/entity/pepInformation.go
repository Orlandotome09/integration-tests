package entity

import "time"

type PepInformation struct {
	DocumentNumber string    `json:"document_number"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	Institution    string    `json:"institution"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Source         string    `json:"source"`
	CreatedAt      time.Time `json:"created_at"`
}

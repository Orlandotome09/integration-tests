package entity

import "github.com/google/uuid"

type Director struct {
	DirectorID uuid.UUID `json:"director_id"`
	Role       string    `json:"role"`
	Person     Person    `json:"person"`
}

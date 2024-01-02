package contracts

import "github.com/google/uuid"

type FindResponse struct {
	ProfileID uuid.UUID       `json:"profile_id"`
	States    []StateResponse `json:"states"`
}

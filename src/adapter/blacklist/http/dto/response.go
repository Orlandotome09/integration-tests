package dto

import "time"

type BlacklistResponse struct {
	PartnerId      string          `json:"partnerId"`
	DocumentNumber string          `json:"documentNumber"`
	Justifications []Justification `json:"justifications"`
	CreatedAt      time.Time       `json:"createdAt,omitempty"`
	UpdatedAt      time.Time       `json:"updatedAt,omitempty"`
	Version        int             `json:"version"`
	Status         BlacklistStatus `json:"status,omitempty"`
}

type BlacklistStatus = string

type Justification struct {
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	Justification string    `json:"justification"`
}

const (
	BlacklistActive   BlacklistStatus = "ACTIVE"
	BlacklistInactive BlacklistStatus = "INACTIVE"
)

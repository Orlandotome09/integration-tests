package entity

import (
	"github.com/google/uuid"
)

type Address struct {
	AddressID    *uuid.UUID `json:"address_id"`
	ProfileID    *uuid.UUID `json:"profile_id"`
	Type         string     `json:"type"`
	ZipCode      string     `json:"zip_code"`
	Street       string     `json:"street"`
	Number       string     `json:"number"`
	Complement   string     `json:"complement"`
	Neighborhood string     `json:"neighborhood"`
	City         string     `json:"city"`
	StateCode    string     `json:"state_code"`
	CountryCode  string     `json:"country_code"`
}

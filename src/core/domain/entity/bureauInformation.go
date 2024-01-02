package entity

import "time"

type BureauInformation struct {
	Name        string     `json:"name"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
}

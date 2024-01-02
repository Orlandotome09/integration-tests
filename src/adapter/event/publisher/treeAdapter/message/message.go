package message

import "time"

type TreeAdapterMessage struct {
	ProfileID          string    `json:"profile_id" binding:"required"`
	Partner            string    `json:"partner" binding:"required"`
	ProfileCreatedDate string    `json:"profile_created_date" binding:"required"` //ddmmyyyy
	Negotiation        string    `json:"negotiation" binding:"required"`
	Person             Person    `json:"person" binding:"required"`
	Date               time.Time `json:"date" binding:"required"`
}

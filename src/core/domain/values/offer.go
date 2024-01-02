package values

import "time"

type Offer struct {
	Type      string    `validate:"required"`
	Product   string    `validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

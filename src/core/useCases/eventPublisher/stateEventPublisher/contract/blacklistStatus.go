package contract

import "time"

type BlacklistStatus struct {
	Status        string        `json:"status,omitempty"`
	Justification Justification `json:"justification,omitempty"`
}

type Justification struct {
	AddedAt  time.Time `json:"added_at,omitempty"`
	Author   string    `json:"author,omitempty"`
	Comments []string  `json:"comments,omitempty"`
}

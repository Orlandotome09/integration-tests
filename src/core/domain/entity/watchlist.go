package entity

import "encoding/json"

type Watchlist struct {
	Sources  []string        `json:"sources"`
	Metadata json.RawMessage `json:"metadata"`
}

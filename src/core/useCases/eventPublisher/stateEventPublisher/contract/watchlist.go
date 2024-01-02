package contract

type Watchlist struct {
	Sources  []string    `json:"sources"`
	Metadata interface{} `json:"metadata"`
}

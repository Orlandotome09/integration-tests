package dto

type WatchlistResponse struct {
	Title   string   `json:"title"`
	Name    string   `json:"name"`
	Link    string   `json:"link"`
	Watch   string   `json:"watch"`
	Other   string   `json:"other"`
	Entries []string `json:"entries"`
	Sources []string `json:"sources"`
}

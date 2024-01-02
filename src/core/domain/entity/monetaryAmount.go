package entity

type MonetaryAmount struct {
	Amount   float64      `json:"amount"`
	Currency CurrencyCode `json:"currency"`
}

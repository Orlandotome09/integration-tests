package contracts

import "github.com/google/uuid"

type AccountResponse struct {
	AccountID     uuid.UUID `json:"account_id"`
	ProfileID     uuid.UUID `json:"profile_id"`
	BankCode      string    `json:"bank_code"`
	AgencyNumber  string    `json:"agency_number"`
	AgencyDigit   string    `json:"agency_digit"`
	AccountNumber string    `json:"account_number"`
	AccountDigit  string    `json:"account_digit"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorsResponse struct {
	Error []string `json:"error"`
}

package entity

import "github.com/google/uuid"

type Account struct {
	AccountID     *uuid.UUID
	ProfileID     *uuid.UUID
	LegacyID      string
	BankCode      string
	AgencyNumber  string
	AgencyDigit   string
	AccountNumber string
	AccountDigit  string
}
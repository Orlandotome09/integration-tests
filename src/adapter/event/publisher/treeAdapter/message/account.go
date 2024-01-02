package message

type Account struct {
	Number       string `json:"number"`
	AgencyNumber string `json:"agency_number"`
	BankCode     string `json:"bank_code"`
}

type Accounts []Account
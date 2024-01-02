package contracts

type CreateAccountRequest struct {
	ProfileID string `json:"profile_id"`
	BankCode  string `json:"bank_code"`
}

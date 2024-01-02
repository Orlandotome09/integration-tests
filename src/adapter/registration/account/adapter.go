package account

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account/http/contracts"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type accountAdapter struct {
	accountClient accountClient.AccountClient
}

func NewAccountAdapter(accountClient accountClient.AccountClient) interfaces.AccountAdapter {
	return &accountAdapter{
		accountClient: accountClient,
	}
}

func (ref *accountAdapter) GetByID(accountID string) (*entity.Account, error) {
	resp, err := ref.accountClient.Get(accountID)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	account := &entity.Account{ProfileID: &resp.ProfileID}

	return account, nil
}

func (ref *accountAdapter) FindByProfileID(profileID string) ([]entity.Account, error) {
	accountResponse, err := ref.accountClient.Find(profileID)
	if err != nil {
		return nil, err
	}

	var accounts []entity.Account
	for _, response := range accountResponse {
		account := entity.Account{
			AccountNumber: response.AccountNumber,
			AccountDigit:  response.AccountDigit,
			AgencyNumber:  response.AgencyNumber,
			AgencyDigit:   response.AgencyDigit,
			BankCode:      response.BankCode,
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (ref *accountAdapter) CreateInternal(entityId string, bankCode string) (*entity.Account, error) {
	request := &contracts.CreateAccountRequest{ProfileID: entityId, BankCode: bankCode}

	resp, err := ref.accountClient.CreateInternal(request)
	if err != nil {
		return nil, err
	}

	account := &entity.Account{ProfileID: &resp.ProfileID}

	return account, nil
}

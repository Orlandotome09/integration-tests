package accountTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type AccountTranslator interface {
	Translate(accounts []entity.Account) message.Accounts
}

type accountTranslator struct{}

func NewAccountTranslator() AccountTranslator {
	return &accountTranslator{}
}

func (ref *accountTranslator) Translate(accounts []entity.Account) message.Accounts {
	messageAccounts := make([]message.Account, len(accounts))
	for i, account := range accounts {
		messageAccounts[i] = message.Account{
			Number:       account.AccountNumber + account.AccountDigit,
			AgencyNumber: account.AgencyNumber + account.AgencyDigit,
			BankCode:     account.BankCode,
		}
	}
	return messageAccounts
}

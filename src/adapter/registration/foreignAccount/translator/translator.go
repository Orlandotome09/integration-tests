package foreignAccountTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ForeignAccountTranslator interface {
	Translate(response contracts.ForeignAccountResponse) entity.ForeignAccount
}

type foreignAccountTranslator struct{}

func New() ForeignAccountTranslator {
	return &foreignAccountTranslator{}
}

func (ref *foreignAccountTranslator) Translate(response contracts.ForeignAccountResponse) entity.ForeignAccount {
	return entity.ForeignAccount{
		ForeignAccountID: response.ForeignAccountID,
		ProfileID:        response.ProfileID,
	}
}

package foreignAccount

import (
	foreignAccountClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http"
	foreignAccountTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type foreignAccountAdapter struct {
	client     foreignAccountClient.ForeignAccountClient
	translator foreignAccountTranslator.ForeignAccountTranslator
}

func NewForeignAccountAdapter(client foreignAccountClient.ForeignAccountClient,
	translator foreignAccountTranslator.ForeignAccountTranslator) interfaces.ForeignAccountAdapter {
	return &foreignAccountAdapter{
		client:     client,
		translator: translator,
	}
}

func (ref *foreignAccountAdapter) Get(foreignAccountID string) (*entity.ForeignAccount, error) {

	resp, err := ref.client.Get(foreignAccountID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, nil
	}

	foreignAccount := ref.translator.Translate(*resp)

	return &foreignAccount, nil
}

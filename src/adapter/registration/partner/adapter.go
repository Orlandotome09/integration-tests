package partner

import (
	partnerClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http"
	partnerTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

var (
	PartnerActive = string(values2.PartnerStatusActive)
)

type partnerAdapter struct {
	client     partnerClient.PartnerClient
	translator partnerTranslator.PartnerTranslator
}

func NewPartnerAdapter(client partnerClient.PartnerClient,
	translator partnerTranslator.PartnerTranslator) interfaces.PartnerAdapter {
	return &partnerAdapter{
		client:     client,
		translator: translator,
	}
}

func (ref *partnerAdapter) GetActive(partnerID string) (*entity.Partner, error) {
	response, err := ref.client.Get(partnerID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if response == nil {
		return nil, values2.NewErrorNotFound("partner")
	}

	if response.Status != PartnerActive {
		return nil, values2.NewErrorValidation("Partner is not active")
	}

	partner := ref.translator.Translate(*response)

	return &partner, nil
}

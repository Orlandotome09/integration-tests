package partnerTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

var (
	StatusActive          = "ACTIVE"
	SegregationByPartner  = "BY_PARTNER"
	SegregationByMerchant = "BY_MERCHANT"
)

type PartnerTranslator interface {
	Translate(response contracts.PartnerResponse) entity.Partner
}

type partnerTranslator struct{}

func New() PartnerTranslator {
	return &partnerTranslator{}
}

func (ref *partnerTranslator) Translate(response contracts.PartnerResponse) entity.Partner {

	return entity.Partner{
		PartnerID:      response.PartnerID,
		DocumentNumber: response.DocumentNumber,
		Name:           response.Name,
		Status:         translateStatus(response.Status),
		LogoImageUrl:   response.LogoImageUrl,
		Config: &entity.PartnerConfig{
			CustomerSegregationType: translateSegregationType(response.Config.CustomerSegregationType),
			UseCallbackV2:           &response.Config.UseCallbackV2,
		},
	}
}

func translateStatus(status string) values2.PartnerStatus {
	switch status {
	case StatusActive:
		return values2.PartnerStatusActive
	default:
		return ""
	}
}

func translateSegregationType(segregationType string) values2.SegregationType {
	switch segregationType {
	case SegregationByPartner:
		return values2.SegregationTypeByPartner
	case SegregationByMerchant:
		return values2.SegregationTypeByMerchant
	default:
		return ""
	}
}

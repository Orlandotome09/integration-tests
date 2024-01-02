package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type OfferBase struct {
	OfferType string `json:"offer_type"`
	Product   string `json:"product" validate:"required"`
}

func (ref OfferBase) ToDomain() *values.Offer {
	offer := &values.Offer{
		Type:    ref.OfferType,
		Product: ref.Product,
	}
	return offer
}

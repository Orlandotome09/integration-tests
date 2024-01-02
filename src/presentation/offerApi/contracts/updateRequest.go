package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type UpdateOfferRequest struct {
	UpdateOfferRequestURI
	UpdateOfferRequestJSON
}

type UpdateOfferRequestURI struct {
	Type string `uri:"offer_type" binding:"required"`
}

type UpdateOfferRequestJSON struct {
	Product string `json:"product" binding:"required"`
}

func (ref UpdateOfferRequest) ToDomain() *values.Offer {
	offer := &values.Offer{
		Type:    ref.Type,
		Product: ref.Product,
	}
	return offer
}

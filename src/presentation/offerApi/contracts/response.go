package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"time"
)

type OfferResponse struct {
	OfferBase
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ref OfferResponse) FromDomain(offer values.Offer) OfferResponse {
	ref.OfferBase.OfferType = offer.Type
	ref.OfferBase.Product = offer.Product
	ref.CreatedAt = offer.CreatedAt
	ref.UpdatedAt = offer.UpdatedAt
	return ref
}

type ListOfferResponse []OfferBase

func (ref ListOfferResponse) FromDomain(offers []values.Offer) ListOfferResponse {
	responses := make([]OfferResponse, len(offers))
	for i, offer := range offers {
		response := OfferResponse{}.FromDomain(offer)
		responses[i] = response
	}
	return ref
}

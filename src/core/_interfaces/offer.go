package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type OfferRepository interface {
	Create(offer values.Offer) (createdOffer *values.Offer, err error)
	Get(offerType string) (offer *values.Offer, err error)
	Update(offer values.Offer) (updatedOffer *values.Offer, err error)
	Delete(offerType string) error
	List() (offers []values.Offer, err error)
}

type OfferService interface {
	Create(offer values.Offer) (createdOffer *values.Offer, err error)
	Get(offerType string) (offer *values.Offer, err error)
	Update(offer values.Offer) (updatedOffer *values.Offer, err error)
	Delete(offerType string) error
	List() (offers []values.Offer, err error)
}

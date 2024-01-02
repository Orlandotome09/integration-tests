package model

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type Offer struct {
	OfferType string `gorm:"primaryKey"`
	Product   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ref Offer) ToDomain() *values.Offer {
	return &values.Offer{
		Type:      ref.OfferType,
		Product:   ref.Product,
		CreatedAt: ref.CreatedAt,
		UpdatedAt: ref.UpdatedAt,
	}
}

func (ref Offer) FromDomain(offer values.Offer) *Offer {
	ref.OfferType = offer.Type
	ref.Product = offer.Product
	return &ref
}

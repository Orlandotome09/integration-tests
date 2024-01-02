package contracts

import (
	 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type Company struct {
	LegalName                 string          `json:"legal_name,omitempty"`
	BusinessName              string          `json:"business_name,omitempty"`
	TaxPayerIdentification    string          `json:"tax_payer_identification,omitempty"`
	CompanyRegistrationNumber string          `json:"company_registration_number,omitempty"`
	DateOfIncorporation       *values.Date    `json:"date_of_incorporation,omitempty"`
	PlaceOfIncorporation      string          `json:"place_of_incorporation,omitempty"`
	ShareCapital              *MonetaryAmount `json:"share_capital,omitempty"`
	License                   string          `json:"license,omitempty"`
	Website                   string          `json:"website,omitempty"`
	GoodsDelivery             *GoodsDelivery  `json:"goods_delivery,omitempty"`
	Assets                    float64         `json:"assets,omitempty"`
	AssetsCurrency            string          `json:"assets_currency,omitempty"`
	AnnualIncome              float64         `json:"annual_income,omitempty"`
	AnnualIncomeCurrency      string          `json:"annual_income_currency,omitempty"`
	CountryCode               string          `json:"country_code,omitempty"`
	LegalNature               string          `json:"legal_nature,omitempty"`
}

func NewCompanyFromDomain(dom *entity.Company) *Company {
	if dom == nil {
		return nil
	}
	return &Company{
		LegalName:                 dom.LegalName,
		BusinessName:              dom.BusinessName,
		TaxPayerIdentification:    dom.TaxPayerIdentification,
		CompanyRegistrationNumber: dom.CompanyRegistrationNumber,
		DateOfIncorporation:       values.NewDateFromTime(dom.DateOfIncorporation),
		PlaceOfIncorporation:      dom.PlaceOfIncorporation,
		ShareCapital:              (&MonetaryAmount{}).FromDomain(dom.ShareCapital),
		License:                   dom.License,
		Website:                   dom.Website,
		GoodsDelivery:             NewGoodsDeliveryFromDomain(dom.GoodsDelivery),
		LegalNature:               dom.LegalNature,
		Assets:                    dom.Assets,
		AssetsCurrency:            dom.AssetsCurrency,
		AnnualIncome:              dom.AnnualIncome,
		AnnualIncomeCurrency:      dom.AnnualIncomeCurrency,
		CountryCode:               dom.CountryCode,
	}
}

func (ref *Company) ToDomain() *entity.Company {
	if ref == nil {
		return nil
	}

	return &entity.Company{
		LegalName:                 ref.LegalName,
		BusinessName:              ref.BusinessName,
		TaxPayerIdentification:    ref.TaxPayerIdentification,
		CompanyRegistrationNumber: ref.CompanyRegistrationNumber,
		DateOfIncorporation:       ref.DateOfIncorporation.ToTime(),
		PlaceOfIncorporation:      ref.PlaceOfIncorporation,
		ShareCapital:              ref.ShareCapital.ToDomain(),
		License:                   ref.License,
		Website:                   ref.Website,
		GoodsDelivery:             ref.GoodsDelivery.ToDomain(),
		Assets:                    ref.Assets,
		AssetsCurrency:            ref.AssetsCurrency,
		AnnualIncome:              ref.AnnualIncome,
		AnnualIncomeCurrency:      ref.AnnualIncomeCurrency,
		CountryCode:               ref.CountryCode,
		LegalNature:               ref.LegalNature,
	}
}

type MonetaryAmount struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func (ref *MonetaryAmount) ToDomain() *entity.MonetaryAmount {
	if ref == nil {
		return nil
	}
	return &entity.MonetaryAmount{
		Amount:   ref.Amount,
		Currency: entity.CurrencyCode(ref.Currency),
	}
}

func (ref *MonetaryAmount) FromDomain(monetary *entity.MonetaryAmount) *MonetaryAmount {
	if monetary == nil {
		return nil
	}
	ref.Amount = monetary.Amount
	ref.Currency = string(monetary.Currency)
	return ref
}

type GoodsDelivery struct {
	AverageDeliveryDays   int    `json:"average_delivery_days"`
	ShippingMethods       string `json:"shipping_methods"`
	Insurance             bool   `json:"insurance"`
	TrackingCodeAvailable bool   `json:"tracking_code_available"`
}

func (ref *GoodsDelivery) ToDomain() *entity.GoodsDelivery {
	if ref == nil {
		return nil
	}

	return &entity.GoodsDelivery{
		AverageDeliveryDays:   ref.AverageDeliveryDays,
		ShippingMethods:       entity.ShippingMethod(ref.ShippingMethods),
		Insurance:             ref.Insurance,
		TrackingCodeAvailable: ref.TrackingCodeAvailable,
	}
}

func NewGoodsDeliveryFromDomain(delivery *entity.GoodsDelivery) *GoodsDelivery {
	if delivery == nil {
		return nil
	}

	return &GoodsDelivery{
		AverageDeliveryDays:   delivery.AverageDeliveryDays,
		ShippingMethods:       string(delivery.ShippingMethods),
		Insurance:             delivery.Insurance,
		TrackingCodeAvailable: delivery.TrackingCodeAvailable,
	}
}

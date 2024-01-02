package contract

import "time"

type Company struct {
	LegalName                 string          `json:"legal_name,omitempty"`
	BusinessName              string          `json:"business_name,omitempty"`
	TaxPayerIdentification    string          `json:"tax_payer_identification,omitempty"`
	CompanyRegistrationNumber string          `json:"company_registration_number,omitempty"`
	DateOfIncorporation       *time.Time      `json:"date_of_incorporation,omitempty"`
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

type MonetaryAmount struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type GoodsDelivery struct {
	AverageDeliveryDays   int    `json:"average_delivery_days"`
	ShippingMethods       string `json:"shipping_methods"`
	Insurance             bool   `json:"insurance"`
	TrackingCodeAvailable bool   `json:"tracking_code_available"`
}

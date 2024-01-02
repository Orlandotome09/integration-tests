package contracts

import "time"

type LegalRepresentativeResponse struct {
	LegalRepresentativeID string     `json:"legal_representative_id"`
	ProfileID             string     `json:"profile_id"`
	PartnerID             string     `json:"partner_id"`
	FullName              string     `json:"full_name"`
	DocumentNumber        string     `json:"document_number"`
	Email                 string     `json:"email"`
	Phone                 string     `json:"phone"`
	Nationality           string     `json:"nationality"`
	BirthDate             string     `json:"birth_date"`
	OfferType             string     `json:"offer_type"`
	Pep                   *bool      `json:"pep"`
	ForeignTaxResidency   bool       `json:"foreign_tax_residency"`
	CountryOfTaxResidency string     `json:"country_of_tax_residency"`
	ExpirationDate        *time.Time `json:"expiration_date,omitempty"`
}

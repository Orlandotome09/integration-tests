package entity

import "time"

type Individual struct {
	FirstName             string             `json:"first_name,omitempty"`
	LastName              string             `json:"last_name,omitempty"`
	DateOfBirth           *time.Time         `json:"date_of_birth,omitempty"`
	DateOfBirthInputted   *time.Time         `json:"date_of_birth_inputted,omitempty"`
	Phones                []Phone            `json:"phones,omitempty"`
	BureauInformation     *BureauInformation `json:"bureau_information,omitempty"`
	Income                *float64           `json:"income,omitempty"`
	IncomeCurrency        string             `json:"income_currency,omitempty"`
	Assets                *float64           `json:"assets,omitempty"`
	AssetsCurrency        string             `json:"assets_currency,omitempty"`
	Nationality           string             `json:"nationality,omitempty"`
	Pep                   *bool              `json:"pep,omitempty"`
	UsPerson              *bool              `json:"us_person,omitempty"`
	Occupation            *string            `json:"occupation,omitempty"`
	ForeignTaxResidency   bool               `json:"foreign_tax_residency"`
	CountryOfTaxResidency string             `json:"country_of_tax_residency,omitempty"`
	Contacts              *string            `json:"contacts,omitempty"`
}

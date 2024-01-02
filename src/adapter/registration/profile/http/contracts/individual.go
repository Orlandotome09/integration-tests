package contracts

import (
	 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"time"
)

type Individual struct {
	FirstName             string             `json:"first_name"`
	LastName              string             `json:"last_name"`
	DateOfBirth           *time.Time         `json:"date_of_birth"`
	DateOfBirthInputted   *time.Time         `json:"date_of_birth_inputted"`
	Phones                Phones             `json:"phones,omitempty"`
	BureauInformation     *BureauInformation `json:"bureau_information,omitempty"`
	Income                *float64           `json:"income,omitempty"`
	IncomeCurrency        string             `json:"income_currency,omitempty"`
	Assets                *float64           `json:"assets,omitempty"`
	AssetsCurrency        string             `json:"assets_currency,omitempty"`
	Nationality           string             `json:"nationality,omitempty"`
	Email                 string             `json:"email,omitempty"`
	Pep                   *bool              `json:"pep,omitempty"`
	UsPerson              *bool              `json:"us_person,omitempty"`
	Occupation            *string            `json:"occupation,omitempty"`
	ForeignTaxResidency   bool               `json:"foreign_tax_residency,omitempty"`
	CountryOfTaxResidency string             `json:"country_of_tax_residency,omitempty"`
}

func NewIndividualFromDomain(individual *entity.Individual) *Individual {
	if individual == nil {
		return nil
	}

	return &Individual{
		FirstName: individual.FirstName,
		LastName:  individual.LastName,
		// TODO Isolate enrichment data
		DateOfBirth:           individual.DateOfBirthInputted,
		DateOfBirthInputted:   individual.DateOfBirthInputted,
		BureauInformation:     NewBureauInformationFromDomain(individual.BureauInformation),
		Phones:                NewPhonesFromDomain(individual.Phones),
		Income:                individual.Income,
		Assets:                individual.Assets,
		Nationality:           individual.Nationality,
		Pep:                   individual.Pep,
		UsPerson:              individual.UsPerson,
		Occupation:            individual.Occupation,
		ForeignTaxResidency:   individual.ForeignTaxResidency,
		CountryOfTaxResidency: individual.CountryOfTaxResidency,
		IncomeCurrency:        individual.IncomeCurrency,
		AssetsCurrency:        individual.AssetsCurrency,
	}
}

func (ref *Individual) ToDomain() *entity.Individual {
	if ref == nil {
		return nil
	}

	return &entity.Individual{
		FirstName:             ref.FirstName,
		LastName:              ref.LastName,
		DateOfBirth:           ref.DateOfBirth,
		Phones:                ref.Phones.ToDomain(),
		DateOfBirthInputted:   ref.DateOfBirthInputted,
		BureauInformation:     ref.BureauInformation.ToDomain(),
		Income:                ref.Income,
		IncomeCurrency:        ref.IncomeCurrency,
		Assets:                ref.Assets,
		AssetsCurrency:        ref.AssetsCurrency,
		Nationality:           ref.Nationality,
		Pep:                   ref.Pep,
		UsPerson:              ref.UsPerson,
		Occupation:            ref.Occupation,
		ForeignTaxResidency:   ref.ForeignTaxResidency,
		CountryOfTaxResidency: ref.CountryOfTaxResidency,
	}
}

type BureauInformation struct {
	Name        string     `json:"name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

func (ref *BureauInformation) ToDomain() *entity.BureauInformation {
	if ref == nil {
		return nil
	}
	return &entity.BureauInformation{
		Name:        ref.Name,
		DateOfBirth: ref.DateOfBirth,
	}
}

func NewBureauInformationFromDomain(dom *entity.BureauInformation) *BureauInformation {
	if dom == nil {
		return nil
	}

	return &BureauInformation{
		Name:        dom.Name,
		DateOfBirth: dom.DateOfBirth,
	}
}

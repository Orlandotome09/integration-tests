package companyTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/timeTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

var (
	EmptyTranslatedDateOfIncorporation  = "00000000"
	EmptyTranslatedPlaceOfIncorporation = ""
)

type CompanyTranslator interface {
	TranslateLegalName(company *entity.Company) string
	TranslateDateOfIncorporation(company *entity.Company) string
	TranslatePlaceOfIncorporation(company *entity.Company) string
	TranslateIncome(company *entity.Company) float64
}

type companyTranslator struct {
	timeTranslator timeTranslator.TimeTranslator
}

func NewCompanyTranslator(timeTranslator timeTranslator.TimeTranslator) CompanyTranslator {
	return &companyTranslator{
		timeTranslator: timeTranslator,
	}
}

func (ref *companyTranslator) TranslateLegalName(company *entity.Company) string {
	if company == nil {
		return ""
	}
	return company.LegalName
}

func (ref *companyTranslator) TranslateDateOfIncorporation(company *entity.Company) string {
	if company == nil || company.DateOfIncorporation == nil {
		return EmptyTranslatedDateOfIncorporation
	}
	return ref.timeTranslator.TranslateTime(*company.DateOfIncorporation)
}

func (ref *companyTranslator) TranslatePlaceOfIncorporation(company *entity.Company) string {
	if company == nil || company.PlaceOfIncorporation == "" {
		return EmptyTranslatedPlaceOfIncorporation
	}
	return company.PlaceOfIncorporation
}

func (ref *companyTranslator) TranslateIncome(company *entity.Company) float64 {
	if company == nil {
		return 0
	}
	return company.AnnualIncome / 12
}

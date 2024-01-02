package companyTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockCompanyTranslator struct {
	CompanyTranslator
	mock.Mock
}

func (ref *MockCompanyTranslator) TranslateLegalName(company *entity.Company) string {
	args := ref.Called(company)
	return args.Get(0).(string)
}

func (ref *MockCompanyTranslator) TranslateDateOfIncorporation(company *entity.Company) string {
	args := ref.Called(company)
	return args.Get(0).(string)
}

func (ref *MockCompanyTranslator) TranslatePlaceOfIncorporation(company *entity.Company) string {
	args := ref.Called(company)
	return args.Get(0).(string)
}

func (ref *MockCompanyTranslator) TranslateIncome(company *entity.Company) float64 {
	args := ref.Called(company)
	return args.Get(0).(float64)
}

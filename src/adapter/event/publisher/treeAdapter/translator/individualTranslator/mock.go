package individualTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockIndividualTranslator struct {
	IndividualTranslator
	mock.Mock
}

func (ref *MockIndividualTranslator) TranslateName(individual *entity.Individual) string {
	args := ref.Called(individual)
	return args.Get(0).(string)
}

func (ref *MockIndividualTranslator) TranslateDateOfBirth(individual *entity.Individual) string {
	args := ref.Called(individual)
	return args.Get(0).(string)
}

func (ref *MockIndividualTranslator) TranslateNationality(individual *entity.Individual) string {
	args := ref.Called(individual)
	return args.Get(0).(string)
}

func (ref *MockIndividualTranslator) TranslatePep(individual *entity.Individual) bool {
	args := ref.Called(individual)
	return args.Get(0).(bool)
}

func (ref *MockIndividualTranslator) TranslateIncome(individual *entity.Individual) float64 {
	args := ref.Called(individual)
	return args.Get(0).(float64)
}

func (ref *MockIndividualTranslator) TranslateAssets(individual *entity.Individual) float64 {
	args := ref.Called(individual)
	return args.Get(0).(float64)
}

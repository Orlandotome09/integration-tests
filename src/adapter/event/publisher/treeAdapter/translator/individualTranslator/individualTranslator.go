package individualTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/timeTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

var (
	EmptyDateOfBirth = "00000000"
	EmptyNationality = ""
)

type IndividualTranslator interface {
	TranslateName(individual *entity.Individual) string
	TranslateDateOfBirth(individual *entity.Individual) string
	TranslateNationality(individual *entity.Individual) string
	TranslatePep(individual *entity.Individual) bool
	TranslateIncome(individual *entity.Individual) float64
	TranslateAssets(individual *entity.Individual) float64
}

type individualTranslator struct {
	timeTranslator timeTranslator.TimeTranslator
}

func NewIndividualTranslator(timeTranslator timeTranslator.TimeTranslator) IndividualTranslator {
	return &individualTranslator{
		timeTranslator: timeTranslator,
	}
}

func (ref *individualTranslator) TranslateName(individual *entity.Individual) string {
	if individual == nil {
		return ""
	}

	if individual.LastName != "" {
		return individual.FirstName + " " + individual.LastName
	}
	return individual.FirstName
}

func (ref *individualTranslator) TranslateDateOfBirth(individual *entity.Individual) string {
	if individual == nil {
		return EmptyDateOfBirth
	}
	if individual.DateOfBirthInputted != nil {
		return ref.timeTranslator.TranslateTime(*individual.DateOfBirthInputted)
	}
	if individual.DateOfBirth != nil {
		return ref.timeTranslator.TranslateTime(*individual.DateOfBirth)
	}
	return EmptyDateOfBirth
}

func (ref *individualTranslator) TranslateNationality(individual *entity.Individual) string {
	if individual == nil || individual.Nationality == "" {
		return EmptyNationality
	}
	return individual.Nationality
}

func (ref *individualTranslator) TranslatePep(individual *entity.Individual) bool {
	if individual == nil || individual.Pep == nil {
		return false
	}
	return *individual.Pep
}

func (ref *individualTranslator) TranslateIncome(individual *entity.Individual) float64 {
	if individual == nil || individual.Income == nil {
		return 0
	}
	return *individual.Income
}

func (ref *individualTranslator) TranslateAssets(individual *entity.Individual) float64 {
	if individual == nil || individual.Assets == nil {
		return 0
	}
	return *individual.Assets
}

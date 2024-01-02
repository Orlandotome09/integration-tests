package person

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

const (
	defaultMinimalAge = 18
)

type isUnderAgeAnalyzer struct {
	person     entity.Person
	minimumAge *int
}

func NewIsUnderAgeAnalyzer(person entity.Person, minimumAge *int) entity.Rule {
	return &isUnderAgeAnalyzer{
		person:     person,
		minimumAge: minimumAge,
	}
}

func (ref *isUnderAgeAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	customerIsUnderAge := entity.NewRuleResultV2(values.RuleSetIsUnderAge, values.RuleNameCustomerIsUnderAge)

	minimumAgeValue := getMinimalAgeValue(ref.minimumAge)

	if hasMinimumAge(ref.person, minimumAgeValue) {
		customerIsUnderAge.SetResult(values.ResultStatusApproved)
		return []entity.RuleResultV2{*customerIsUnderAge}, nil
	}

	customerIsUnderAge.SetResult(values.ResultStatusRejected)
	return []entity.RuleResultV2{*customerIsUnderAge}, nil
}

func (ref *isUnderAgeAnalyzer) Name() values.RuleSet {
	return values.RuleSetIsUnderAge
}

func getMinimalAgeValue(minimalAge *int) int {
	if minimalAge != nil {
		return *minimalAge
	}

	return defaultMinimalAge
}

func hasMinimumAge(person entity.Person, minimalAge int) bool {
	if person.PersonType == values.PersonTypeIndividual {
		if person.Individual == nil || person.Individual.DateOfBirth == nil {
			return false
		}
		personAge := time.Now().Year() - person.Individual.DateOfBirth.Year()

		return personAge >= minimalAge
	}

	return true
}

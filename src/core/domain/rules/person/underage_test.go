package person

import (
	"testing"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
)

func TestIsUnderAgeAnalyzer_should_approve_when_individual_is_over_default_minimum_age_18(t *testing.T) {
	now := time.Now()
	dateOfBirth := now.AddDate(-18, 0, 0)
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{
			DateOfBirth: &dateOfBirth,
		},
	}

	rulesResult, err := NewIsUnderAgeAnalyzer(person, nil).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
}

func TestIsUnderAgeAnalyzer_should_approve_when_individual_is_over_age_16(t *testing.T) {
	now := time.Now()
	dateOfBirth := now.AddDate(-18, 0, 0)
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{
			DateOfBirth: &dateOfBirth,
		},
	}
	minimumAge := 16

	rulesResult, err := NewIsUnderAgeAnalyzer(person, &minimumAge).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
}

func TestIsUnderAgeAnalyzer_should_reject_when_individual_is_under_age_18(t *testing.T) {
	now := time.Now()
	dateOfBirth := now.AddDate(-17, 0, 0)
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{
			DateOfBirth: &dateOfBirth,
		},
	}

	rulesResult, err := NewIsUnderAgeAnalyzer(person, nil).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusRejected, rulesResult[0].Result)
}

func TestIsUnderAgeAnalyzer_should_reject_when_individual_is_under_age_16(t *testing.T) {
	now := time.Now()
	dateOfBirth := now.AddDate(-15, 0, 0)
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{
			DateOfBirth: &dateOfBirth,
		},
	}
	minimumAge := 16

	rulesResult, err := NewIsUnderAgeAnalyzer(person, &minimumAge).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusRejected, rulesResult[0].Result)
}

func TestIsUnderAgeAnalyzer_company_is_always_over_age(t *testing.T) {
	person := entity.Person{
		PersonType: values.PersonTypeCompany,
	}
	minimalAge := 999

	rulesResult, err := NewIsUnderAgeAnalyzer(person, &minimalAge).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
}

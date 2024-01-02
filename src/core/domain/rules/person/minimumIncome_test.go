package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinimumIncomeAnalyzer_Analyze_Over(t *testing.T) {
	personIncome := 2000.00
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{
			Income: &personIncome,
		},
	}

	rulesResult, err := NewMinimumIncomeAnalyzer(person).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
	assert.False(t, rulesResult[0].Pending)
}

func TestMinimumIncomeAnalyzer_Analyze_Equal(t *testing.T) {
	personIncome := 1000.00
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{
			Income: &personIncome,
		},
	}

	rulesResult, err := NewMinimumIncomeAnalyzer(person).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
	assert.False(t, rulesResult[0].Pending)
}

func TestMinimumIncomeAnalyzer_Analyze_Below(t *testing.T) {
	personIncome := 999.00
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{
			Income: &personIncome,
		},
	}

	rulesResult, err := NewMinimumIncomeAnalyzer(person).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, rulesResult[0].Result, values.ResultStatusAnalysing)
	assert.True(t, rulesResult[0].Pending)
}

func TestMinimumIncomeAnalyzer_PersonTypeIndividual_Invalid(t *testing.T) {
	person := entity.Person{
		PersonType: values.PersonTypeCompany,
	}

	expectedError := "Minimum income rule is only applied to individual"
	rulesResult, err := NewMinimumIncomeAnalyzer(person).Analyze()

	assert.Nil(t, rulesResult)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestMinimumIncomeAnalyzer_Individual_Nil(t *testing.T) {
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
	}

	expectedError := "Individual is required for minimum income rule"
	rulesResult, err := NewMinimumIncomeAnalyzer(person).Analyze()

	assert.Nil(t, rulesResult)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestMinimumIncomeAnalyzer_Income_NotDeclared(t *testing.T) {
	person := entity.Person{
		PersonType: values.PersonTypeIndividual,
		Individual: &entity.Individual{},
	}

	rulesResult, err := NewMinimumIncomeAnalyzer(person).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, rulesResult[0].Result, values.ResultStatusAnalysing)
	assert.True(t, rulesResult[0].Pending)
}

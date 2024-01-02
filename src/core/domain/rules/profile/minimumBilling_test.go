package profile

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinimumBillingAnalyzer_Analyze_Over(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
			Company: &entity.Company{
				AnnualIncome: 2000.00,
			},
		},
	}

	rulesResult, err := NewMinimumBillingAnalyzer(profile).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
	assert.False(t, rulesResult[0].Pending)
}

func TestMinimumBillingAnalyzer_Analyze_Equal(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
			Company: &entity.Company{
				AnnualIncome: 1000.00,
			},
		},
	}

	rulesResult, err := NewMinimumBillingAnalyzer(profile).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
	assert.False(t, rulesResult[0].Pending)
}

func TestMinimumBillingAnalyzer_Analyze_Below(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
			Company: &entity.Company{
				AnnualIncome: 999.00,
			},
		},
	}

	rulesResult, err := NewMinimumBillingAnalyzer(profile).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, rulesResult[0].Result, values.ResultStatusAnalysing)
	assert.True(t, rulesResult[0].Pending)
}

func TestMinimumBillingAnalyzer_PersonTypeIndividual_Invalid(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			PersonType: values.PersonTypeIndividual,
		},
	}

	expectedError := "Minimum billing rule is only applied to company"
	rulesResult, err := NewMinimumBillingAnalyzer(profile).Analyze()

	assert.Nil(t, rulesResult)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestMinimumBillingAnalyzer_Company_Nil(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
		},
	}

	expectedError := "Company is required for minimum billing rule"
	rulesResult, err := NewMinimumBillingAnalyzer(profile).Analyze()

	assert.Nil(t, rulesResult)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestMinimumBillingAnalyzer_AnnualIncome_NotDeclared(t *testing.T) {
	profile := entity.Profile{
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
			Company:    &entity.Company{},
		},
	}

	rulesResult, err := NewMinimumBillingAnalyzer(profile).Analyze()

	assert.Nil(t, err)
	assert.Equal(t, rulesResult[0].Result, values.ResultStatusAnalysing)
	assert.True(t, rulesResult[0].Pending)
}

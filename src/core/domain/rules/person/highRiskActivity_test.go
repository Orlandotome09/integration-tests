package person

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

func TestAnalyze_shouldThrowErrorForProfileNotCompany(t *testing.T) {
	profileID := uuid.New()
	person := entity.Person{ProfileID: profileID, DocumentNumber: "123", PersonType: values.PersonTypeIndividual}
	economicActivityService := &mocks.EconomicActivityService{}
	riskyActivityRule := NewHighRiskActivityAnalyzer(person, economicActivityService)

	var expected []entity.RuleResultV2 = nil
	received, err := riskyActivityRule.Analyze()

	assert.NotNil(t, err)
	assert.Equal(t, expected, received)
}

func TestAnalyze_shouldNotFoundCompanyInBureau(t *testing.T) {
	profileID := uuid.New()
	person := entity.Person{ProfileID: profileID, DocumentNumber: "123", OfferType: "xx1", PartnerID: "p23", PersonType: values.PersonTypeCompany}
	economicActivityService := &mocks.EconomicActivityService{}
	riskyActivityRule := NewHighRiskActivityAnalyzer(person, economicActivityService)

	metadata, _ := json.Marshal(fmt.Sprintf("company (%s) not found in bureau", person.DocumentNumber))
	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetActivityRisk,
			RuleName: values.RuleNameHighRiskActivity,
			Result:   values.ResultStatusRejected,
			Metadata: metadata,
		},
	}

	received, err := riskyActivityRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestAnalyze_shouldFindActivityOfHighRisk(t *testing.T) {
	person := entity.Person{
		DocumentNumber: "123",
		PersonType:     values.PersonTypeCompany,
		EnrichedInformation: &entity.EnrichedInformation{
			EnrichedCompany: entity.EnrichedCompany{
				EconomicActivity: "any code",
			},
		},
	}
	economicActivityService := &mocks.EconomicActivityService{}
	riskyActivityRule := NewHighRiskActivityAnalyzer(person, economicActivityService)

	economicActivityCode := person.EnrichedInformation.EconomicActivity
	economicActivity := &entity.EconomicActivity{RiskValue: true}

	economicActivityService.On("Get", economicActivityCode).Return(economicActivity, true, nil)

	metadata, _ := json.Marshal(fmt.Sprintf("economic activity (%s) is high risk", economicActivityCode))
	activityCode := ActivityCode{Code: economicActivityCode}
	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetActivityRisk,
			RuleName: values.RuleNameHighRiskActivity,
			Result:   values.ResultStatusAnalysing,
			Pending:  true,
			Metadata: metadata,
			Problems: []entity.Problem{{Code: values.ProblemCodeEconomicalActivityRiskHigh, Detail: activityCode}},
		},
	}

	received, err := riskyActivityRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestAnalyze_shouldFindNotRiskyActivity(t *testing.T) {
	person := entity.Person{
		DocumentNumber: "123",
		PersonType:     values.PersonTypeCompany,
		EnrichedInformation: &entity.EnrichedInformation{
			EnrichedCompany: entity.EnrichedCompany{
				EconomicActivity: "any code",
			},
		},
	}

	economicActivityService := &mocks.EconomicActivityService{}
	riskyActivityRule := NewHighRiskActivityAnalyzer(person, economicActivityService)

	businessActivity := person.EnrichedInformation.EconomicActivity
	economicalActivity := &entity.EconomicActivity{RiskValue: false}

	economicActivityService.On("Get", businessActivity).Return(economicalActivity, true, nil)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetActivityRisk,
			RuleName: values.RuleNameHighRiskActivity,
			Result:   values.ResultStatusApproved,
		},
	}

	received, err := riskyActivityRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestAnalyze_shouldNotFindActivityInRepository(t *testing.T) {
	person := entity.Person{
		DocumentNumber: "123",
		PersonType:     values.PersonTypeCompany,
		EnrichedInformation: &entity.EnrichedInformation{
			EnrichedCompany: entity.EnrichedCompany{
				EconomicActivity: "any code",
			},
		},
	}

	economicActivityService := &mocks.EconomicActivityService{}
	riskyActivityRule := NewHighRiskActivityAnalyzer(person, economicActivityService)

	economicActivityCode := person.EnrichedInformation.EconomicActivity
	economicalActivity := &entity.EconomicActivity{RiskValue: true}

	economicActivityService.On("Get", economicActivityCode).Return(economicalActivity, false, nil)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetActivityRisk,
			RuleName: values.RuleNameHighRiskActivity,
			Result:   values.ResultStatusApproved,
		},
	}

	received, err := riskyActivityRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestName(t *testing.T) {
	riskyActivityRule := NewHighRiskActivityAnalyzer(entity.Person{}, nil)

	expected := values.RuleSetActivityRisk
	received := riskyActivityRule.Name()

	assert.Equal(t, expected, received)
}

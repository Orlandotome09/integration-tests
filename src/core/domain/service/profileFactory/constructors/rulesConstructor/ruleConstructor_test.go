package rulesConstructor

import (
	profileRulesFactory "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_When_Should_Get_Profile_Rules(t *testing.T) {

	profileID := uuid.New()

	ruleSetConfig := entity.RuleSetConfig{
		BureauParams: &entity.BureauParams{},
	}

	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						StepNumber:  1,
						RulesConfig: &ruleSetConfig,
					},
				},
			},
		},
	}

	profileWrapper := entity.ProfileWrapper{
		Profile: profile,
	}

	rulesFactory := &profileRulesFactory.ProfileRulesFactory{}

	rules := []entity.Rule{}

	rulesFactory.On("GetRules", &ruleSetConfig, &profile).Return(rules, nil)

	constructor := profileRulesConstructor{profileRulesFactory: rulesFactory}

	err := constructor.Assemble(&profileWrapper)

	expected := profile
	expected.ValidationSteps = []entity.RuleValidatorStep{
		{
			StepNumber:      1,
			SkipForApproval: false,
			Rules:           rules,
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, &expected.ValidationSteps, &profileWrapper.Profile.ValidationSteps)
	mock.AssertExpectationsForObjects(t, rulesFactory)

}

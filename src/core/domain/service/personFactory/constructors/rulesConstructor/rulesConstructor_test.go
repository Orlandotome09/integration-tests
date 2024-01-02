package personRulesConstructor

import (
	personRulesFactory "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_Should_Populate_Rules(t *testing.T) {
	personRulesFactoryInstance := &personRulesFactory.PersonRulesFactory{}
	constructor := personRulesConstructor{
		personRulesFactory: personRulesFactoryInstance,
	}

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					StepNumber:      1,
					SkipForApproval: false,
					RulesConfig: &entity.RuleSetConfig{
						BureauParams: &entity.BureauParams{},
					},
				},
			},
		},
	}
	personWrapper := entity.PersonWrapper{
		Person: person,
	}
	var rules []entity.Rule

	personRulesFactoryInstance.On("GetRules", *person.CadastralValidationConfig.ValidationSteps[0].RulesConfig, person).Return(rules, nil)

	err := constructor.Assemble(&personWrapper)

	expected := []entity.RuleValidatorStep{
		{
			StepNumber:      1,
			SkipForApproval: false,
			Rules:           rules,
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.ValidationSteps)
	mock.AssertExpectationsForObjects(t, personRulesFactoryInstance)
}

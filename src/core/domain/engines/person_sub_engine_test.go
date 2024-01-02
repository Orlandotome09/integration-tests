package engines

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type personSubEngineParams struct {
	ruleValidator    *mocks.RuleValidator
	personFactory    *mocks.PersonFactory
	personRepository *mocks.PersonRepository
}

func preparePersonSubEngine() (personSubEngineParams, _interfaces.SubEngine) {
	ruleValidator := &mocks.RuleValidator{}
	personFactory := &mocks.PersonFactory{}
	personRepository := &mocks.PersonRepository{}

	subEngine := NewPersonSubEngine(ruleValidator, personFactory, personRepository)
	return personSubEngineParams{
		ruleValidator:    ruleValidator,
		personFactory:    personFactory,
		personRepository: personRepository,
	}, subEngine
}

func TestPersonSubEngine_Prepare(t *testing.T) {
	params, subEngine := preparePersonSubEngine()

	entityID := uuid.New()

	person := entity.Person{
		EntityID:                  entityID,
		RoleType:                  values.RoleTypeShareholder,
		OfferType:                 "SOME OFFER",
		PersonType:                values.PersonTypeIndividual,
		CadastralValidationConfig: &entity.CadastralValidationConfig{},
	}

	params.personFactory.On("Build", person).Return(&person, nil)
	params.ruleValidator.On("SetRules", person.ValidationSteps).Return()
	params.personRepository.On("Save", person).Return(&person, nil)

	err := subEngine.Prepare(person, "SOME OFFER")

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, params.ruleValidator, params.personFactory, params.personRepository)
}

func TestPersonSubEngine_Validate(t *testing.T) {
	params, subEngine := preparePersonSubEngine()

	state := entity.State{}
	override := entity.Overrides{}
	noCache := false
	entityID := uuid.New()
	engineName := "PROFILE"

	err := errors.New("no error")
	params.ruleValidator.On("Validate", state, override, noCache, entityID, engineName).Return(&state, err)

	expectedErr := err
	expected := &state
	received, receivedErr := subEngine.Validate(state, override, noCache, entityID, engineName)

	assert.Equal(t, expectedErr, receivedErr)
	assert.Equal(t, expected, received)
}

func TestPersonSubEngine_GetName(t *testing.T) {
	_, subEngine := preparePersonSubEngine()

	expected := values.EngineNamePerson
	received := subEngine.GetName()

	assert.Equal(t, expected, received)
}

func TestPersonSubEngine_NewInstance(t *testing.T) {
	params, subEngine := preparePersonSubEngine()

	newRuleValidator := &mocks.RuleValidator{}

	params.ruleValidator.On("NewInstance").Return(newRuleValidator)

	expected := &personSubEngine{
		ruleValidator:    newRuleValidator,
		personFactory:    params.personFactory,
		personRepository: params.personRepository,
	}

	received := subEngine.NewInstance()

	assert.Equal(t, expected, received)
}

package engines

import (
	"errors"
	"testing"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type profileEngineParams struct {
	profileValidator  *mocks.RuleValidator
	profileAdapter    *mocks.ProfileAdapter
	posProcessor      *mocks.PosProcessor
	profileFactory    *mocks.ProfileFactory
	profileRepository *mocks.ComplianceProfileRepository
}

func prepareProfileEngine(profile entity.Profile) (profileEngineParams, interfaces.Engine) {
	profileValidator := &mocks.RuleValidator{}
	profileService := &mocks.ProfileAdapter{}
	posProcessor := &mocks.PosProcessor{}
	profileFactory := &mocks.ProfileFactory{}
	profileRepository := &mocks.ComplianceProfileRepository{}
	engine := &profileEngine{
		profile:           profile,
		validator:         profileValidator,
		posProcessor:      posProcessor,
		profileAdapter:    profileService,
		profileFactory:    profileFactory,
		profileRepository: profileRepository,
	}

	return profileEngineParams{
		profileValidator:  profileValidator,
		profileAdapter:    profileService,
		posProcessor:      posProcessor,
		profileFactory:    profileFactory,
		profileRepository: profileRepository,
	}, engine
}

func TestProfileEngine_Prepare(t *testing.T) {
	t.Run("should prepare successfully", func(t *testing.T) {
		params, engine := prepareProfileEngine(entity.Profile{})

		entityID := uuid.New()
		profile := &entity.Profile{
			Person: entity.Person{
				PersonType:                values.PersonTypeIndividual,
				OfferType:                 "any",
				RoleType:                  values.RoleTypeCustomer,
				CadastralValidationConfig: &entity.CadastralValidationConfig{},
			},
		}

		params.profileFactory.On("Build", entityID).Return(profile, nil)
		params.profileValidator.On("SetRules", profile.ValidationSteps)
		params.profileRepository.On("Save", *profile).Return(profile, nil)

		err := engine.Prepare(entityID)

		assert.Nil(t, err)
		mock.AssertExpectationsForObjects(t, params.profileFactory, params.profileValidator, params.profileRepository)
	})

	t.Run("should not prepare when profile not found", func(t *testing.T) {
		params, engine := prepareProfileEngine(entity.Profile{})
		entityID := uuid.New()

		params.profileFactory.On("Build", entityID).Return(nil, nil)

		err := engine.Prepare(entityID)

		assert.ErrorContains(t, err, "profile not found")
		mock.AssertExpectationsForObjects(t, params.profileFactory)
	})

	t.Run("should not prepare when not found cadastral validation config", func(t *testing.T) {
		params, engine := prepareProfileEngine(entity.Profile{})
		entityID := uuid.New()
		profile := &entity.Profile{
			Person: entity.Person{
				PersonType:                values.PersonTypeIndividual,
				OfferType:                 "any",
				RoleType:                  values.RoleTypeCustomer,
				CadastralValidationConfig: nil,
			},
		}

		params.profileFactory.On("Build", entityID).Return(profile, nil)

		err := engine.Prepare(entityID)

		assert.ErrorContains(t, err, "cadastral validation config not found for profile")
		mock.AssertExpectationsForObjects(t, params.profileFactory, params.profileValidator, params.profileRepository)
		params.profileValidator.AssertNumberOfCalls(t, "SetRules", 0)
		params.profileRepository.AssertNumberOfCalls(t, "Save", 0)
	})
}

func TestProfileEngine_Validate(t *testing.T) {
	params, engine := prepareProfileEngine(entity.Profile{})

	state := entity.State{}
	override := entity.Overrides{}
	noCache := false
	entityID := uuid.New()
	engineName := "PROFILE"
	err := errors.New("no error")

	params.profileValidator.On("Validate", state, override, noCache, entityID, engineName).Return(&state, err)

	expectedErr := err
	expected := &state
	received, receivedErr := engine.Validate(state, override, noCache, entityID, engineName)

	assert.Equal(t, expectedErr, receivedErr)
	assert.Equal(t, expected, received)
}

func TestProfileEngine_PosProcessing(t *testing.T) {
	t.Run("should create internal account", func(t *testing.T) {
		previousState := &entity.State{Result: values.ResultStatusAnalysing}
		state := &entity.State{Result: values.ResultStatusApproved}
		entityID := uuid.New()
		profile := entity.Profile{
			Person: entity.Person{
				OfferType:                 "xx",
				PersonType:                "xx",
				RoleType:                  values.RoleTypeCustomer,
				CadastralValidationConfig: &entity.CadastralValidationConfig{ValidationSteps: []entity.ValidationStep{}, ProductConfig: &entity.ProductConfig{CreateBexsAccount: true}},
			},
		}

		params, engine := prepareProfileEngine(profile)

		params.posProcessor.On("SendToTreeAdapter", profile, *profile.CadastralValidationConfig.ProductConfig).Return(nil)
		params.posProcessor.On("SendToLimit", profile, *state, *profile.CadastralValidationConfig.ProductConfig).Return(nil)
		params.posProcessor.On("CreateInternalAccount", entityID).Return(nil)

		err := engine.PosProcessing(previousState, state, entityID)

		assert.Nil(t, err)
		mock.AssertExpectationsForObjects(t, params.posProcessor)
	})

	t.Run("should not create internal account", func(t *testing.T) {
		previousState := &entity.State{Result: values.ResultStatusAnalysing}
		newState := &entity.State{Result: values.ResultStatusApproved}
		entityID := uuid.New()

		profile := entity.Profile{
			Person: entity.Person{
				OfferType:                 "xx",
				PersonType:                "xx",
				RoleType:                  values.RoleTypeCustomer,
				CadastralValidationConfig: &entity.CadastralValidationConfig{ValidationSteps: []entity.ValidationStep{}, ProductConfig: &entity.ProductConfig{CreateBexsAccount: false}},
			},
		}

		params, engine := prepareProfileEngine(profile)

		params.posProcessor.On("SendToTreeAdapter", profile, *profile.CadastralValidationConfig.ProductConfig).Return(nil)
		params.posProcessor.On("SendToLimit", profile, *newState, *profile.CadastralValidationConfig.ProductConfig).Return(nil)

		err := engine.PosProcessing(previousState, newState, entityID)

		assert.Nil(t, err)
		params.posProcessor.AssertNumberOfCalls(t, "CreateInternalAccount", 0)
		mock.AssertExpectationsForObjects(t, params.posProcessor)
	})
}

func TestProfileEngine_GetName(t *testing.T) {
	_, engine := prepareProfileEngine(entity.Profile{})

	expected := values.EngineNameProfile
	received := engine.GetName()

	assert.Equal(t, expected, received)
}

func TestProfileEngine_NewInstance(t *testing.T) {
	params, engine := prepareProfileEngine(entity.Profile{})

	newProfileValidator := &mocks.RuleValidator{}
	params.profileValidator.On("NewInstance").Return(newProfileValidator)

	expected := &profileEngine{
		validator:         newProfileValidator,
		profileAdapter:    params.profileAdapter,
		posProcessor:      params.posProcessor,
		profileFactory:    params.profileFactory,
		profileRepository: params.profileRepository,
	}

	received := engine.NewInstance()

	assert.Equal(t, expected, received)
}

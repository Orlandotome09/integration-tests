package complianceAnalyzer

import (
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type complianceAnalyzer struct {
	overrideRepository interfaces.OverrideRepository
	stateManager       interfaces.StateManager
}

type validationInput struct {
	complianceValidator interfaces.ComplianceValidator
	state               *entity.State
	overrides           entity.Overrides
	EntityID            uuid.UUID
	noCache             bool
}

func NewComplianceAnalyzer(overrideRepository interfaces.OverrideRepository,
	stateManager interfaces.StateManager) interfaces.ComplianceAnalyzer {
	return &complianceAnalyzer{
		overrideRepository: overrideRepository,
		stateManager:       stateManager,
	}
}

func (ref *complianceAnalyzer) RunComplianceAnalysis(
	complianceValidator interfaces.ComplianceValidator,
	entityID uuid.UUID,
	requestDate time.Time,
	executionTime time.Time,
	cacheValue bool,
) (*entity.State, error) {

	overrides, err := ref.overrideRepository.FindByEntityID(entityID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	state, shouldIgnore, err := ref.stateManager.GetOrCreateState(requestDate, entityID, complianceValidator.GetName())
	if err != nil {
		return nil, errors.WithStack(err)
	} else if shouldIgnore {
		return state, nil
	}

	previousState := *state
	input := validationInput{
		complianceValidator: complianceValidator,
		state:               state,
		overrides:           overrides,
		EntityID:            entityID,
		noCache:             cacheValue,
	}

	newState, err := ref.execute(input)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := ref.stateManager.UpdateState(newState, requestDate, executionTime); err != nil {
		return nil, errors.WithStack(err)
	}

	if err = complianceValidator.PosProcessing(&previousState, newState, entityID); err != nil {
		return nil, errors.WithStack(err)
	}

	return newState, nil
}

func (ref *complianceAnalyzer) execute(input validationInput) (*entity.State, error) {
	shouldSkipValidations := (input.overrides.HasBlocked() && input.state.Result == values.ResultStatusBlocked) || input.overrides.HasInactive()

	if shouldSkipValidations {
		logrus.Infof("[complianceAnalyzer]Skipping validations since entity is Blocked or Inactive. Entity ID: %v. Override Blocked: %v. Current Result: %v", input.EntityID, input.overrides.HasBlocked(), input.state.Result)
		return input.state, nil
	}

	stateResult, err := input.complianceValidator.Validate(*input.state, input.overrides, input.noCache, input.EntityID,
		input.complianceValidator.GetName())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logrus.WithField("state", stateResult).
		Infof("[complianceAnalyzer] New state for Entity %s", input.EntityID)
	return stateResult, nil

}

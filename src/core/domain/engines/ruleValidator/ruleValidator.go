package ruleValidator

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	RuleValidatorClassName = "[RuleValidator]"
	validateMethod         = "[validate]"
)

type RuleValidator struct {
	validationSteps []entity.RuleValidatorStep
	StateMachine    interfaces.StateMachine
}

func New(stateMachine interfaces.StateMachine) interfaces.RuleValidator {
	return &RuleValidator{
		StateMachine: stateMachine,
	}
}

func (ref *RuleValidator) Validate(state entity.State, override entity.Overrides, noCache bool, entityID uuid.UUID, engineName string) (*entity.State, error) {

	if len(ref.validationSteps) == 0 {
		state.Result = values.ResultStatusApproved
		return &state, nil
	}

	for _, validationStep := range ref.validationSteps {

		resultState, err := validate(state, override, noCache, validationStep, entityID, engineName)
		if err != nil {
			return nil, err
		}

		state = *ref.StateMachine.CalculateState(*resultState)

		//Continue only if all steps for the category were approved so far
		result := state.GetStepResult(validationStep.StepNumber)
		if result != values.ResultStatusApproved {
			break
		}

	}

	return &state, nil
}

func (ref *RuleValidator) SetRules(rules []entity.RuleValidatorStep) {

	ref.validationSteps = rules
}

func (ref *RuleValidator) NewInstance() interfaces.RuleValidator {
	newRuleValidator := &RuleValidator{
		validationSteps: nil,
		StateMachine:    ref.StateMachine,
	}
	return newRuleValidator
}

func validate(state entity.State, overrides entity.Overrides,
	noCache bool, validationStep entity.RuleValidatorStep, entityID uuid.UUID, engineName string) (*entity.State, error) {

	logrus.WithFields(logrus.Fields{
		"class": RuleValidatorClassName,
		"msg_detail": fmt.Sprintf("Rules to run: %v Previous state: %v",
			getAllRulesNames(validationStep.Rules), state),
	}).Info(fmt.Sprintf("Starting validations with Engine %v for entity(id) %v", engineName, entityID))

	newRulesResult := []entity.RuleResultV2{}

	for _, rule := range validationStep.Rules {
		logrus.WithFields(logrus.Fields{
			"class": RuleValidatorClassName,
		}).Info(fmt.Sprintf("Running rule %v for entity: %v", rule.Name(), entityID))

		results, err := rule.Analyze()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"class": RuleValidatorClassName,
				"error": err,
			}).Error(fmt.Sprintf("Rule %v failed for entity: %v", rule.Name(), entityID))
			return nil, errors.WithStack(err)
		}
		newRulesResult = append(newRulesResult, results...)

		logrus.WithFields(logrus.Fields{
			"class":      RuleValidatorClassName,
			"msg_detail": fmt.Sprintf("New results are: %v", results),
		}).Info(fmt.Sprintf("Finished rule %v for entity: %v", rule.Name(), entityID))

	}

	newRulesResult = applyOverridesOnResults(overrides, newRulesResult)

	state = *updateStepResult(state, validationStep, newRulesResult)

	return &state, nil
}

func getAllRulesNames(rules []entity.Rule) (names []values.RuleSet) {
	for _, rule := range rules {
		names = append(names, rule.Name())
	}

	return
}

func applyOverridesOnResults(overrides entity.Overrides, results []entity.RuleResultV2) []entity.RuleResultV2 {

	if overrides == nil {
		return results
	}
	if results == nil {
		return nil
	}

	updated := []entity.RuleResultV2{}

	for _, result := range results {
		override, exists := overrides.FindByRuleSetAndName(result.RuleSet, result.RuleName)

		if exists {
			metadata := map[string]string{
				"comments": override.Comments,
				"author":   override.Author,
			}

			metadataJson, _ := json.Marshal(metadata)

			result.Result = override.Result
			result.Pending = false
			result.Metadata = metadataJson
		}
		updated = append(updated, result)
	}

	return updated
}

func updateStepResult(state entity.State, validationStep entity.RuleValidatorStep, rulesResult []entity.RuleResultV2) *entity.State {

	if state.ValidationStepsResults == nil {
		state.ValidationStepsResults = make([]entity.ValidationStepResult, 0)
	}

	for index, step := range state.ValidationStepsResults {
		if step.StepNumber == validationStep.StepNumber {
			step.RuleResults = rulesResult
			state.ValidationStepsResults[index] = step
			return &state
		}
	}

	levelResult := &entity.ValidationStepResult{
		StepNumber:      validationStep.StepNumber,
		Result:          "",
		SkipForApproval: validationStep.SkipForApproval,
		RuleResults:     rulesResult,
	}

	state.ValidationStepsResults = append(state.ValidationStepsResults, *levelResult)
	return &state
}

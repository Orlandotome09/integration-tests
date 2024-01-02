package statemachine

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type stateMachine struct {
}

func NewStateMachine() _interfaces.StateMachine {
	return &stateMachine{}
}

func (ref *stateMachine) CalculateState(state entity.State) *entity.State {

	var result values.Result

	for index, step := range state.ValidationStepsResults {
		step.Result = ref.calculateStepResult(step)
		state.ValidationStepsResults[index] = step
	}

	for _, step := range state.ValidationStepsResults {
		// Saves the result of the first Step Non-Approved or Approved with Skip for Approval flag false
		if step.Result != values.ResultStatusApproved || !step.SkipForApproval {
			result = step.Result
			break
		}
	}

	state.Result = result

	pending, ruleNames := ref.aggregateByRules(state.ValidationStepsResults)
	state.Pending = pending
	state.RuleNames = ruleNames

	return &state
}

func (ref *stateMachine) aggregateByRules(validationStepsResults []entity.ValidationStepResult) (pending bool, ruleNames []values.RuleName) {
	for _, step := range validationStepsResults {
		for _, rule := range step.RuleResults {
			pending = pending || rule.Pending
			ruleNames = append(ruleNames, rule.RuleName)
		}
	}
	return
}

func (ref *stateMachine) calculateStepResult(stepResult entity.ValidationStepResult) values.Result {

	worstStatusObtained := values.ResultStatusApproved

	for _, element := range stepResult.RuleResults {
		if element.Result.IsItWorstThan(&worstStatusObtained) {
			worstStatusObtained = element.Result
		}
	}

	return worstStatusObtained
}

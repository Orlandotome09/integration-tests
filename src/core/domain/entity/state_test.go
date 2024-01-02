package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStateGetStepResult(t *testing.T) {
	state := State{ValidationStepsResults: ValidationStepsResult{ValidationStepResult{StepNumber: 1, Result: values.ResultStatusApproved}}}

	expected := values.ResultStatusApproved
	actual := state.GetStepResult(1)

	assert.Equal(t, expected, actual)
}

func TestStateGetStepResultNotFound(t *testing.T) {
	state := State{ValidationStepsResults: ValidationStepsResult{ValidationStepResult{StepNumber: 1, Result: values.ResultStatusApproved}}}

	expected := values.Result("")
	actual := state.GetStepResult(2)

	assert.Equal(t, expected, actual)
}

func TestStateIsApproved(t *testing.T) {
	state := State{Result: values.ResultStatusApproved}

	expected := true
	actual := state.IsApproved()

	assert.Equal(t, expected, actual)
}

func TestStateIsNotApproved(t *testing.T) {
	state := State{Result: values.ResultStatusRejected}

	expected := false
	actual := state.IsApproved()

	assert.Equal(t, expected, actual)
}

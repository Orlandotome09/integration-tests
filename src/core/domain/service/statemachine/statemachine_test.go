package statemachine

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfileStateMachine_CalculateState_Rejected(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusRejected, false, false).build()
	state = NewStateMachine().CalculateState(*state)

	if state.Result != values.ResultStatusRejected {
		t.Error("Invalid output for Profile State Machine")
	}

	if state.ValidationStepsResults[0].Result != values.ResultStatusRejected {
		t.Error("Invalid output for Profile State Machine")
	}

	if state.Pending {
		t.Error("Invalid output for Profile State Machine")
	}

	if state.RuleNames[0] != "TEST" {
		t.Error("Invalid output for Profile State Machine")
	}

	if len(state.RuleNames) != 1 {
		t.Error("Invalid output for Profile State Machine")
	}
}

func TestProfileStateMachine_CalculateState_Analysing(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusAnalysing, true, false).build()
	state = NewStateMachine().CalculateState(*state)

	if state.Result != values.ResultStatusAnalysing {
		t.Error("Invalid output for Profile State Machine")
	}

	if !state.Pending {
		t.Error("Invalid output for Profile State Machine")
	}

	if state.RuleNames[0] != "TEST" {
		t.Error("Invalid output for Profile State Machine")
	}

	if len(state.RuleNames) != 1 {
		t.Error("Invalid output for Profile State Machine")
	}
}

func TestProfileStateMachine_CalculateState_Approved(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusApproved, false, false).build()
	state = NewStateMachine().CalculateState(*state)

	if state.Result != values.ResultStatusApproved {
		t.Error("Invalid output for Profile State Machine")
	}

	if state.ValidationStepsResults[0].Result != values.ResultStatusApproved {
		t.Error("Invalid output for Profile State Machine")
	}

	if state.Pending {
		t.Error("Invalid output for Profile State Machine")
	}

	if state.RuleNames[0] != "TEST" {
		t.Error("Invalid output for Profile State Machine")
	}

	if len(state.RuleNames) != 1 {
		t.Error("Invalid output for Profile State Machine")
	}
}

func TestProfileStateMachine_CalculateState_Blocked(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusBlocked, false, false).build()
	state = NewStateMachine().CalculateState(*state)

	if state.Result != values.ResultStatusBlocked {
		t.Error("Invalid output for Profile State Machine")
	}
}

func TestProfileStateMachine_CalculateState_Inactive(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusInactive, false, false).build()
	state = NewStateMachine().CalculateState(*state)

	if state.Result != values.ResultStatusInactive {
		t.Error("Invalid output for Profile State Machine")
	}
}

func TestProfileStateMachine_CalculateState_MultipleSteps_SkipForApprovalFlag(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusApproved, false, true).
		withAdditionalStep(values.ResultStatusRejected, false, false).
		build()
	state = NewStateMachine().CalculateState(*state)

	assert.Equal(t, values.ResultStatusRejected, state.Result)
}

func TestProfileStateMachine_CalculateState_MultipleSteps_No_SkipForApprovalFlag(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusApproved, false, false).
		withAdditionalStep(values.ResultStatusRejected, false, false).
		build()
	state = NewStateMachine().CalculateState(*state)

	assert.Equal(t, values.ResultStatusApproved, state.Result)
}

func TestProfileStateMachine_CalculateState_MultipleSteps_SkipForApprovalFlag_Rejected(t *testing.T) {
	state := newProfileStateBuilder().withResult(values.ResultStatusRejected, false, true).
		withAdditionalStep(values.ResultStatusIgnored, false, false).
		build()
	state = NewStateMachine().CalculateState(*state)

	assert.Equal(t, values.ResultStatusRejected, state.Result)
}

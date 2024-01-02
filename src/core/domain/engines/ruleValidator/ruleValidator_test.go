package ruleValidator

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_updateStepResult_new_step(t *testing.T) {

	state := entity.State{}

	validationStep := entity.RuleValidatorStep{
		StepNumber: 0,
	}

	ruleResults := []entity.RuleResultV2{}

	result := updateStepResult(state, validationStep, ruleResults)

	expected := entity.State{
		ValidationStepsResults: []entity.ValidationStepResult{
			{
				StepNumber:      0,
				SkipForApproval: false,
				RuleResults:     ruleResults,
			},
		},
	}

	assert.Equal(t, *result, expected)

}

func Test_updateStepResult_existing_step(t *testing.T) {

	state := entity.State{
		ValidationStepsResults: []entity.ValidationStepResult{
			{
				StepNumber:      0,
				SkipForApproval: true,
				RuleResults:     nil,
			},
		},
	}

	validationStep := entity.RuleValidatorStep{
		StepNumber: 0,
	}

	ruleResults := []entity.RuleResultV2{}

	result := updateStepResult(state, validationStep, ruleResults)

	expected := entity.State{
		ValidationStepsResults: []entity.ValidationStepResult{
			{
				StepNumber:      0,
				SkipForApproval: true,
				RuleResults:     ruleResults,
			},
		},
	}

	assert.Equal(t, *result, expected)

}

func Test_updateStepResult_append_step(t *testing.T) {

	state := entity.State{
		ValidationStepsResults: []entity.ValidationStepResult{
			{
				StepNumber:      0,
				SkipForApproval: true,
				RuleResults:     nil,
			},
		},
	}

	validationStep := entity.RuleValidatorStep{
		StepNumber: 1,
	}

	ruleResults := []entity.RuleResultV2{}

	result := updateStepResult(state, validationStep, ruleResults)

	expected := entity.State{
		ValidationStepsResults: []entity.ValidationStepResult{
			{
				StepNumber:      0,
				SkipForApproval: true,
				RuleResults:     nil,
			},
			{
				StepNumber:      1,
				SkipForApproval: false,
				RuleResults:     ruleResults,
			},
		},
	}

	assert.Equal(t, *result, expected)

}

func Test_applyOverridesOnResults(t *testing.T) {

	overrides := entity.Overrides{
		{
			RuleSet:  "1",
			RuleName: "X",
			Result:   values.ResultStatusRejected,
			Author:   "SomePerson",
			Comments: "SomeComment",
		},
	}

	ruleResults := []entity.RuleResultV2{
		{
			RuleSet:  "1",
			RuleName: "X",
			Result:   values.ResultStatusApproved,
		},
	}

	result := applyOverridesOnResults(overrides, ruleResults)

	metadataJson, _ := json.Marshal(map[string]string{
		"comments": "SomeComment",
		"author":   "SomePerson",
	})

	expected := []entity.RuleResultV2{
		{
			RuleSet:  "1",
			RuleName: "X",
			Result:   values.ResultStatusRejected,
			Metadata: metadataJson,
		},
	}

	assert.Equal(t, result, expected)

}

func Test_applyOverridesOnResults_NoOverrides(t *testing.T) {

	overrides := entity.Overrides{}

	ruleResults := []entity.RuleResultV2{
		{
			RuleSet:  "1",
			RuleName: "X",
			Result:   values.ResultStatusApproved,
		},
	}

	result := applyOverridesOnResults(overrides, ruleResults)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  "1",
			RuleName: "X",
			Result:   values.ResultStatusApproved,
		},
	}

	assert.Equal(t, result, expected)

}

func Test_validate(t *testing.T) {

	state := entity.State{}
	overrides := entity.Overrides{}
	rule := mocks.Rule{}

	validationStep := entity.RuleValidatorStep{
		StepNumber:      0,
		SkipForApproval: false,
		Rules:           []entity.Rule{&rule},
	}

	ruleResults := []entity.RuleResultV2{
		{
			RuleSet:  "1",
			RuleName: "X",
			Result:   values.ResultStatusApproved,
		},
	}

	rule.On("Analyze").Return(ruleResults, nil)
	rule.On("Name").Return(values.RuleSet("1"))

	stateResult, err := validate(state, overrides, true, validationStep, uuid.New(), "")

	expectedState := entity.State{
		ValidationStepsResults: []entity.ValidationStepResult{
			{
				StepNumber:  0,
				RuleResults: ruleResults,
			},
		},
	}

	assert.Equal(t, *stateResult, expectedState)
	assert.Nil(t, err)
}

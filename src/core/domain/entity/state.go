package entity

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type State struct {
	EntityID               uuid.UUID             `json:"entity_id"`
	EngineName             values.EngineName     `json:"engine_name"`
	Result                 values.Result         `json:"result"`
	ValidationStepsResults ValidationStepsResult `json:"validation_steps_results"`
	RuleNames              []values.RuleName     `json:"rule_names"`
	Pending                bool                  `json:"pending"`
	ExecutionTime          time.Time             `json:"execution_time"`
	RequestDate            time.Time             `json:"request_date"`
	Version                int                   `json:"version"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
}

func (state State) GetStepResult(stepNumber int) values.Result {
	for _, step := range state.ValidationStepsResults {
		if step.StepNumber == stepNumber {
			return step.Result
		}
	}
	return ""
}

func (state State) IsApproved() bool {
	if state.Result == values.ResultStatusApproved {
		return true
	}
	return false
}

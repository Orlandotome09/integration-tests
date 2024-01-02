package model

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Aqui o modelo está conhecendo o elemento do domínio => Remover!
type State struct {
	EntityID               uuid.UUID              `gorm:"primaryKey;type:uuid"`
	EngineName             values.EngineName      `gorm:"index"`
	Result                 string                 `gorm:"index"`
	ValidationStepsResults ValidationStepsResults `gorm:"type:jsonb"`
	RuleNames              pq.StringArray         `gorm:"type:text[]"`
	Pending                bool                   `gorm:"index"`
	RequestDate            time.Time
	ExecutionTime          time.Time
	Version                int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// Set custom table name
func (ref State) TableName() string {
	return "profile_states"
}

func StateFromDomain(state *entity.State) *State {

	validationStepsResults := make([]ValidationStepResult, len(state.ValidationStepsResults))
	for i, step := range state.ValidationStepsResults {
		validationStepsResults[i] = *ValidationStepResultFromDomain(&step)
	}
	if len(validationStepsResults) == 0 {
		validationStepsResults = nil
	}

	ruleNames := make([]string, len(state.RuleNames))
	for idx, ruleName := range state.RuleNames {
		ruleNames[idx] = ruleName.ToString()
	}

	return &State{
		EntityID:               state.EntityID,
		EngineName:             state.EngineName,
		Result:                 string(state.Result),
		ValidationStepsResults: validationStepsResults,
		RuleNames:              ruleNames,
		Pending:                state.Pending,
		RequestDate:            state.RequestDate,
		ExecutionTime:          state.ExecutionTime,
		Version:                state.Version,
		CreatedAt:              state.CreatedAt,
		UpdatedAt:              state.UpdatedAt,
	}
}

func (ref *State) ToDomain() *entity.State {

	validationStepResults := make([]entity.ValidationStepResult, len(ref.ValidationStepsResults))
	for i, step := range ref.ValidationStepsResults {
		validationStepResults[i] = *step.ToDomain()
	}

	ruleNames := make([]values.RuleName, len(ref.RuleNames))
	for idx, ruleName := range ref.RuleNames {
		ruleNames[idx] = values.RuleName(ruleName)
	}

	return &entity.State{
		EntityID:               ref.EntityID,
		EngineName:             ref.EngineName,
		Result:                 values.Result(ref.Result),
		ValidationStepsResults: validationStepResults,
		RuleNames:              ruleNames,
		Pending:                ref.Pending,
		RequestDate:            ref.RequestDate,
		ExecutionTime:          ref.ExecutionTime,
		Version:                ref.Version,
		CreatedAt:              ref.CreatedAt,
		UpdatedAt:              ref.UpdatedAt,
	}
}

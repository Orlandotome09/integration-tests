package entity

import (
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/google/uuid"
)

type Notification struct {
	EntityID      uuid.UUID
	EngineName  values2.EngineName
	Result      values2.Result
	CallbackUrl string
	PartnerID     string
	UseCallbackV2 bool
	ResultDetails []ResultDetail
}

func (ref *Notification) FromState(state *State) (*Notification, error) {
	ref.EntityID = state.EntityID
	ref.EngineName = state.EngineName
	ref.Result = state.Result

	ref.ResultDetails = make([]ResultDetail, len(state.ValidationStepsResults))
	for i, step := range state.ValidationStepsResults {
		ref.ResultDetails[i] = *ResultDetailFromValidationStepResult(step)
	}

	return ref, nil
}

func (ref *Notification) CalculateHash(uuidNamespace uuid.UUID) uuid.UUID {
	content, _ := json.Marshal(ref)
	hash := uuid.NewSHA1(uuidNamespace, content)
	return hash
}

type ResultDetail struct {
	Result    values2.Result
	LimitType string
	LimitInterval string
	ApprovedValue *float64
	Problems      []Problem
}

func ResultDetailFromValidationStepResult(validationStepResult ValidationStepResult) *ResultDetail {
	problems := []Problem{}

	//If not approved, see all RuleSetResults and all RuleResults to append Problems

	for _, ruleResult := range validationStepResult.RuleResults {
		if ruleResult.Result != values2.ResultStatusApproved && ruleResult.Problems != nil {
			for _, problem := range ruleResult.Problems {
				problems = append(problems, problem)
			}
		}
	}

	return &ResultDetail{
		Result:   validationStepResult.Result,
		Problems: problems,
	}
}

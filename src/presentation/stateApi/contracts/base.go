package contracts

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
)

type StateBase struct {
	EntityID      string                 `json:"entity_id"`
	EngineName    string                 `json:"engine_name"`
	Result        values2.Result         `json:"result"`
	RuleSetResult []DetailedRuleResponse `json:"rule_set_result,omitempty"`
}

type DetailedRuleResponse struct {
	StepNumber int                   `json:"step_number"`
	Set        values2.RuleSet       `json:"set"`
	Name       values2.RuleName      `json:"name"`
	Result     values2.Result        `json:"result"`
	Pending    bool                  `json:"pending"`
	Metadata   json.RawMessage       `json:"metadata,omitempty" swaggertype:"object"`
	Tags       []string              `json:"tags,omitempty"`
	Problems   []DetailedRuleProblem `json:"problems,omitempty"`
}

type DetailedRuleProblem struct {
	Code   values2.ProblemCode `json:"code"`
	Detail interface{}         `json:"detail"`
}

func (ref *StateBase) FromDomain(state *entity2.State, onlyPending bool) {

	ref.EntityID = state.EntityID.String()
	ref.EngineName = state.EngineName
	ref.Result = state.Result

	ref.RuleSetResult = DetailedRuleResponse{}.FromDomain(state.ValidationStepsResults, onlyPending)

}

func (ref DetailedRuleResponse) FromDomain(validationStepsResults []entity2.ValidationStepResult, onlyPending bool) (output []DetailedRuleResponse) {

	for _, step := range validationStepsResults {
		for _, ruleResult := range step.RuleResults {

			if !onlyPending || ruleResult.Pending {
				output = append(output, DetailedRuleResponse{
					StepNumber: step.StepNumber,
					Set:        ruleResult.RuleSet,
					Name:       ruleResult.RuleName,
					Result:     ruleResult.Result,
					Pending:    ruleResult.Pending,
					Metadata:   ruleResult.Metadata,
					Tags:       ruleResult.Tags,
					Problems:   DetailedRuleProblem{}.FromDomain(ruleResult.Problems),
				})
			}
		}
	}

	return
}

func (ref DetailedRuleProblem) FromDomain(problems []entity2.Problem) (output []DetailedRuleProblem) {
	for _, problem := range problems {
		output = append(output, DetailedRuleProblem{
			Code:   problem.Code,
			Detail: problem.Detail,
		})
	}
	return
}

func (ref DetailedRuleResponse) FromDomainRuleResult(results []entity2.RuleResultV2, onlyPending bool) (output []DetailedRuleResponse) {

	for _, item := range results {
		if !onlyPending || item.Pending {
			output = append(output, DetailedRuleResponse{
				Set:      item.RuleSet,
				Name:     item.RuleName,
				Result:   item.Result,
				Pending:  item.Pending,
				Metadata: item.Metadata,
				Tags:     item.Tags,
				Problems: DetailedRuleProblem{}.FromDomain(item.Problems),
			})
		}
	}

	return
}

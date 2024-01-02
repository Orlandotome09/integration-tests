package model

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

type ValidationStepsResults []ValidationStepResult

func (ref ValidationStepsResults) Value() (driver.Value, error) {
	if ref == nil {
		return nil, nil
	}
	return json.Marshal(ref)
}

func (ref *ValidationStepsResults) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Failed to unmarshal JSONB value: %+v", value)
	}

	return json.Unmarshal(bytes, ref)
}

type ValidationStepResult struct {
	StepNumber      int            `json:"step_number"`
	Result          string         `json:"result"`
	SkipForApproval bool           `json:"skip_for_approval"`
	RuleResults     []RuleResultV2 `json:"rule_results,omitempty"`
}

func ValidationStepResultFromDomain(domainResult *entity2.ValidationStepResult) *ValidationStepResult {
	if domainResult == nil {
		return nil
	}

	ruleResults := make([]RuleResultV2, len(domainResult.RuleResults))
	for i, ruleResult := range domainResult.RuleResults {
		ruleResults[i] = *RuleResultV2FromDomain(&ruleResult)
	}
	if len(ruleResults) == 0 {
		ruleResults = nil
	}

	return &ValidationStepResult{
		StepNumber:      domainResult.StepNumber,
		Result:          string(domainResult.Result),
		RuleResults:     ruleResults,
		SkipForApproval: domainResult.SkipForApproval,
	}
}

func (ref *ValidationStepResult) ToDomain() *entity2.ValidationStepResult {
	if ref == nil {
		return nil
	}

	var domRuleResults = make([]entity2.RuleResultV2, len(ref.RuleResults))
	for i, ruleResult := range ref.RuleResults {
		domRuleResults[i] = *ruleResult.ToDomain()
	}

	return &entity2.ValidationStepResult{
		StepNumber:      ref.StepNumber,
		Result:          values2.Result(ref.Result),
		RuleResults:     domRuleResults,
		SkipForApproval: ref.SkipForApproval,
	}
}

type RuleResultV2 struct {
	RuleSet  string          `json:"rule_set"`
	RuleName string          `json:"rule_name"`
	Result   string          `json:"result"`
	ExpireAt *time.Time      `json:"expire_at,omitempty"`
	Metadata json.RawMessage `json:"metadata"`
	Pending  bool            `json:"pending"`
	Tags     json.RawMessage `json:"tags"`
	Problems []Problem       `json:"problems,omitempty"`
}

func RuleResultV2FromDomain(dom *entity2.RuleResultV2) *RuleResultV2 {
	if dom == nil {
		return nil
	}

	problems := make([]Problem, len(dom.Problems))
	for i, domProblem := range dom.Problems {
		problems[i] = *ProblemFromDomain(&domProblem)
	}
	if len(problems) == 0 {
		problems = nil
	}

	tags, _ := json.Marshal(dom.Tags)

	return &RuleResultV2{
		RuleSet:  string(dom.RuleSet),
		RuleName: string(dom.RuleName),
		Result:   string(dom.Result),
		ExpireAt: dom.ExpireAt,
		Metadata: dom.Metadata,
		Pending:  dom.Pending,
		Tags:     tags,
		Problems: problems,
	}

}

func (ref *RuleResultV2) ToDomain() *entity2.RuleResultV2 {
	if ref == nil {
		return nil
	}

	domProblems := make([]entity2.Problem, len(ref.Problems))
	for i, problem := range ref.Problems {
		domProblems[i] = *problem.ToDomain()
	}

	var tags []string

	if ref.Tags != nil {
		json.Unmarshal(ref.Tags, &tags)
	}

	return &entity2.RuleResultV2{
		RuleSet:  values2.RuleSet(ref.RuleSet),
		RuleName: values2.RuleName(ref.RuleName),
		Result:   values2.Result(ref.Result),
		ExpireAt: ref.ExpireAt,
		Metadata: ref.Metadata,
		Pending:  ref.Pending,
		Tags:     tags,
		Problems: domProblems,
	}
}

type Problem struct {
	Code   string      `json:"code"`
	Detail interface{} `json:"detail,omitempty"`
}

func ProblemFromDomain(dom *entity2.Problem) *Problem {
	if dom == nil {
		return nil
	}

	return &Problem{
		Code:   string(dom.Code),
		Detail: dom.Detail,
	}
}

func (ref *Problem) ToDomain() *entity2.Problem {
	if ref == nil {
		return nil
	}

	return &entity2.Problem{
		Code:   values2.ProblemCode(ref.Code),
		Detail: ref.Detail,
	}
}

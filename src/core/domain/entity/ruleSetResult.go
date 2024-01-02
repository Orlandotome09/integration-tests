package entity

import (
	values "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"time"
)

type RuleResultV2 struct {
	RuleSet           values.RuleSet  `json:"rule_set"`
	RuleName          values.RuleName `json:"rule_name"`
	Result            values.Result   `json:"result"`
	Pending           bool            `json:"pending"`
	ExpireAt          *time.Time      `json:"expire_at"`
	Metadata          json.RawMessage `json:"metadata"`
	Tags              []string        `json:"tags"`
	Problems          []Problem       `json:"problems"`
	ApprovedDocuments []string        `json:"approved_documents"`
}

func (ruleResult RuleResultV2) IsApproved() bool {
	return ruleResult.Result == values.ResultStatusApproved
}

func NewRuleResultV2(ruleSet values.RuleSet, ruleName values.RuleName) *RuleResultV2 {
	return &RuleResultV2{
		RuleSet:  ruleSet,
		RuleName: ruleName,
		Result:   values.ResultStatusIgnored,
		Pending:  false,
		ExpireAt: nil,
		Metadata: nil,
		Tags:     nil,
		Problems: nil,
	}
}

func (ruleResult *RuleResultV2) SetResult(result values.Result) *RuleResultV2 {
	ruleResult.Result = result
	return ruleResult
}

func (ruleResult *RuleResultV2) SetPending(pending bool) *RuleResultV2 {
	ruleResult.Pending = pending
	return ruleResult
}

func (ruleResult *RuleResultV2) SetExpireAt(expireAt *time.Time) *RuleResultV2 {
	ruleResult.ExpireAt = expireAt
	return ruleResult
}

func (ruleResult *RuleResultV2) SetMetadata(metadata []byte) *RuleResultV2 {
	ruleResult.Metadata = metadata
	return ruleResult
}

func (ruleResult *RuleResultV2) AddTag(tag string) *RuleResultV2 {
	if ruleResult.Tags == nil {
		ruleResult.Tags = []string{}
	}
	ruleResult.Tags = append(ruleResult.Tags, tag)
	return ruleResult
}

func (ruleResult *RuleResultV2) AddProblem(problemCode values.ProblemCode, detail interface{}) *RuleResultV2 {
	if ruleResult.Problems == nil {
		ruleResult.Problems = []Problem{}
	}

	if values.ProblemCodeParser(problemCode) != nil {
		ruleResult.Problems = append(ruleResult.Problems, Problem{
			Code:   values.ProblemCodeInvalidProblemCode,
			Detail: problemCode,
		})
		return ruleResult
	}

	ruleResult.Problems = append(ruleResult.Problems, Problem{
		Code:   problemCode,
		Detail: detail,
	})

	return ruleResult
}

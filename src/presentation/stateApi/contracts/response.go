package contracts

import (
	"strings"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type StateResponse struct {
	StateBase
}

type ComplianceResponse struct {
	EntityID       uuid.UUID        `json:"entity_id"`
	EntityType     string           `json:"entity_type,omitempty"`
	Result         values.Result    `json:"result"`
	DetailedResult []detailedResult `json:"detailed_result,omitempty"`
}

type detailedResult struct {
	Result        string   `json:"result,omitempty"`
	LimitType     string   `json:"limit_type,omitempty"`
	LimitInterval string   `json:"limit_interval,omitempty"`
	ApprovedValue *float64 `json:"approved_value,omitempty"`
	Details       []detail `json:"details,omitempty"`
}

type detail struct {
	Code   string      `json:"code,omitempty"`
	Detail interface{} `json:"detail,omitempty"`
}

type ResyncResponse struct {
	Resynced []string `json:"resynced"`
}

type ReprocessedResponse struct {
	Reprocessed []string `json:"reprocessed"`
}

func (ref *ComplianceResponse) FromDomain(state entity.State) {

	ref.EntityID = state.EntityID
	ref.EntityType = state.EngineName
	ref.Result = state.Result

	ref.DetailedResult = make([]detailedResult, 0)
	details := make([]detail, 0)

	if ref.Result == values.ResultStatusApproved {
		return
	}

	for _, step := range state.ValidationStepsResults {

		if step.Result == values.ResultStatusApproved {
			continue
		}

		for _, ruleResult := range step.RuleResults {

			if ruleResult.Result == values.ResultStatusApproved || ruleResult.Result == values.ResultStatusIgnored {
				continue
			}

			metadata := string(ruleResult.Metadata)
			metadata = strings.ReplaceAll(metadata, "\"", "")
			metadata = strings.ReplaceAll(metadata, "[", "")
			metadata = strings.ReplaceAll(metadata, "]", "")

			for _, prob := range ruleResult.Problems {
				details = append(details, detail{
					Detail: prob.Detail,
					Code:   prob.Code,
				})
			}
		}
		ref.DetailedResult = append(ref.DetailedResult, detailedResult{
			Result:  string(step.Result),
			Details: details,
		})

	}

}

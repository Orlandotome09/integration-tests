package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

const (
	cafRejectedStatus = "REJECTED"
	cafApprovedStatus = "APPROVED"
)

type cafAnalyzer struct {
	enrichedInformation *entity.EnrichedInformation
}

func NewCafAnalyzer(enrichedInformation *entity.EnrichedInformation) entity.Rule {
	return &cafAnalyzer{
		enrichedInformation: enrichedInformation,
	}
}

func (ref *cafAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	cafAnalysis := entity.NewRuleResultV2(values.RuleSetCafAnalysis, values.RuleNameCafAnalysis)

	if ref.enrichedInformation == nil {
		cafAnalysis.SetResult(values.ResultStatusAnalysing).SetPending(true)
		cafAnalysis.AddProblem(values.ProblemCodeNotFoundEnrichedInformation, "")

		return []entity.RuleResultV2{*cafAnalysis}, nil
	}

	cafProvider := ref.enrichedInformation.FilterCafProvider()
	if cafProvider == nil {
		cafAnalysis.SetResult(values.ResultStatusAnalysing).SetPending(true)
		cafAnalysis.AddProblem(values.ProblemCodeNotFoundCafAnalysis, "")

		return []entity.RuleResultV2{*cafAnalysis}, nil
	}

	switch cafProvider.Status {
	case cafRejectedStatus:
		cafAnalysis.SetResult(values.ResultStatusRejected)
		return []entity.RuleResultV2{*cafAnalysis}, nil
	case cafApprovedStatus:
		cafAnalysis.SetResult(values.ResultStatusApproved)
		return []entity.RuleResultV2{*cafAnalysis}, nil
	default:
		cafAnalysis.SetResult(values.ResultStatusAnalysing).SetPending(true)
		cafAnalysis.AddProblem(values.ProblemCodeCafAnalysisPending, []string{cafProvider.RequestID})
		return []entity.RuleResultV2{*cafAnalysis}, nil
	}
}

func (ref *cafAnalyzer) Name() values.RuleSet {
	return values.RuleSetCafAnalysis
}

package shareholding

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

const (
	MinShareholding = 95.0
)

type shareholdingRule struct {
	profile entity.Profile
}

func NewShareholdingRule(profile entity.Profile) interfaces.ShareholdingAnalyzer {
	return &shareholdingRule{
		profile: profile,
	}
}

func (ref *shareholdingRule) Analyze() (*entity.RuleResultV2, *entity.OwnershipStructure, error) {
	shareholdingRuleResult := entity.NewRuleResultV2(values.RuleSetOwnershipStructure, values.RuleNameShareholding)

	minShareholdingSum := ref.getMinShareholding()

	enrichedPercentage := 0.0
	if ref.profile.EnrichedInformation != nil && ref.profile.EnrichedInformation.OwnershipStructure != nil {
		enrichedPercentage = ref.profile.EnrichedInformation.OwnershipStructure.ShareholdingSum
	}

	manuallyFilledPercentage := 0.0
	if ref.profile.OwnershipStructure != nil {
		manuallyFilledPercentage = ref.profile.OwnershipStructure.ShareholdingSum
	}

	if enrichedPercentage >= minShareholdingSum {
		shareholdingRuleResult.SetResult(values.ResultStatusApproved).SetPending(false)
		return shareholdingRuleResult, ref.profile.EnrichedInformation.OwnershipStructure, nil
	}

	if manuallyFilledPercentage >= minShareholdingSum {
		shareholdingRuleResult.SetResult(values.ResultStatusApproved).SetPending(false)
		return shareholdingRuleResult, ref.profile.OwnershipStructure, nil
	}

	result := fmt.Sprintf("Shareholing does not achieve %v minimun required. Shareholding enriched: %v. Shareholding manually filled: %v", minShareholdingSum, enrichedPercentage, manuallyFilledPercentage)
	metadata, err := json.Marshal(result)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	shareholdingRuleResult.SetResult(values.ResultStatusAnalysing).SetPending(true).SetMetadata(metadata).
		AddProblem(values.ProblemCodeShareholdingNotAchieveMinimumRequired, map[string]interface{}{
			"min_shareholding":    minShareholdingSum,
			"enriched_percentage": enrichedPercentage,
			"manually_percentage": manuallyFilledPercentage})
	return shareholdingRuleResult, ref.selectOwnershipStructure(), nil

}

func (ref *shareholdingRule) getMinShareholding() float64 {

	if ref.profile.CadastralValidationConfig == nil {
		return MinShareholding
	}

	for _, step := range ref.profile.CadastralValidationConfig.ValidationSteps {
		if step.RulesConfig == nil {
			continue
		}
		if step.RulesConfig.OwnershipStructureParams != nil && step.RulesConfig.OwnershipStructureParams.MinShareholdingPercentage != nil {
			return *step.RulesConfig.OwnershipStructureParams.MinShareholdingPercentage
		}
	}
	return MinShareholding
}

func (ref *shareholdingRule) selectOwnershipStructure() *entity.OwnershipStructure {
	if ref.profile.OwnershipStructure != nil {
		return ref.profile.OwnershipStructure
	}
	if ref.profile.EnrichedInformation != nil && ref.profile.EnrichedInformation.OwnershipStructure != nil {
		return ref.profile.EnrichedInformation.OwnershipStructure
	}
	return nil
}

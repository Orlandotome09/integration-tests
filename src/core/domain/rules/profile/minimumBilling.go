package profile

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

const (
	minimumBillingDefault = 1000.00
)

type minimumBillingAnalyzer struct {
	ProfileRule
}

func NewMinimumBillingAnalyzer(profile entity.Profile) entity.Rule {
	return &minimumBillingAnalyzer{
		ProfileRule: ProfileRule{
			profile: profile,
		},
	}
}

func (ref *minimumBillingAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	results, err := ref.validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if results != nil {
		return results, err
	}

	minimumBillingResult := entity.NewRuleResultV2(values.RuleSetMinimumBilling, values.RuleNameInsufficientBilling)
	minimumBillingRequired := ref.getMinimumBillingRequired()

	if ref.profile.Company.AnnualIncome >= minimumBillingRequired {
		minimumBillingResult.SetResult(values.ResultStatusApproved).SetPending(false)
		return []entity.RuleResultV2{*minimumBillingResult}, nil
	}

	result := fmt.Sprintf("Company billing is below %v", minimumBillingRequired)
	metadata, err := json.Marshal(result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	minimumBillingResult.SetResult(values.ResultStatusAnalysing).SetPending(true).SetMetadata(metadata).
		AddProblem(values.ProblemCodeCompanyHasInsufficientBilling, map[string]interface{}{
			"minimum_billing_required": minimumBillingRequired,
			"company_billing":          ref.profile.Company.AnnualIncome,
		})

	return []entity.RuleResultV2{*minimumBillingResult}, nil
}

func (ref *minimumBillingAnalyzer) getMinimumBillingRequired() float64 {
	if ref.profile.CadastralValidationConfig == nil {
		return minimumBillingDefault
	}

	for _, step := range ref.profile.CadastralValidationConfig.ValidationSteps {
		if step.RulesConfig == nil {
			continue
		}
		if step.RulesConfig.MinimumBillingParams != nil && step.RulesConfig.MinimumBillingParams.MinimumBilling != nil {
			return *step.RulesConfig.MinimumBillingParams.MinimumBilling
		}
	}

	return minimumBillingDefault
}

func (ref *minimumBillingAnalyzer) validate() ([]entity.RuleResultV2, error) {
	if ref.profile.PersonType != values.PersonTypeCompany {
		return nil, errors.New("Minimum billing rule is only applied to company")
	}

	if ref.profile.Company == nil {
		return nil, errors.New("Company is required for minimum billing rule")
	}

	return nil, nil
}

func (ref *minimumBillingAnalyzer) Name() values.RuleSet {
	return values.RuleSetMinimumBilling
}

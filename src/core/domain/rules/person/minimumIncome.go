package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

const (
	minimumIncomeDefault = 1000.00
)

type minimumIncomeAnalyzer struct {
	person entity.Person
}

func NewMinimumIncomeAnalyzer(person entity.Person) entity.Rule {
	return &minimumIncomeAnalyzer{
		person: person,
	}
}

func (ref *minimumIncomeAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	results, err := ref.validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if results != nil {
		return results, err
	}

	if ref.person.Individual.Income == nil {
		ref.person.Individual.Income = new(float64)
	}

	minimumIncomeResult := entity.NewRuleResultV2(values.RuleSetMinimumIncome, values.RuleNameInsufficientIncome)
	minimumIncomeRequired := ref.getMinimumIncomeRequired()

	if *ref.person.Individual.Income >= minimumIncomeRequired {
		minimumIncomeResult.SetResult(values.ResultStatusApproved).SetPending(false)
		return []entity.RuleResultV2{*minimumIncomeResult}, nil
	}

	result := fmt.Sprintf("Profile income is below %v", minimumIncomeRequired)
	metadata, err := json.Marshal(result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	minimumIncomeResult.SetResult(values.ResultStatusAnalysing).SetPending(true).SetMetadata(metadata).
		AddProblem(values.ProblemCodePersonHasInsufficientMinimumIncome, map[string]interface{}{
			"minimum_income_required": minimumIncomeRequired,
			"person_income":           *ref.person.Individual.Income,
		})

	return []entity.RuleResultV2{*minimumIncomeResult}, nil
}

func (ref *minimumIncomeAnalyzer) getMinimumIncomeRequired() float64 {

	if ref.person.CadastralValidationConfig == nil {
		return minimumIncomeDefault
	}

	for _, step := range ref.person.CadastralValidationConfig.ValidationSteps {
		if step.RulesConfig == nil {
			continue
		}
		if step.RulesConfig.MinimumIncomeParams != nil && step.RulesConfig.MinimumIncomeParams.MinimumIncome != nil {
			return *step.RulesConfig.MinimumIncomeParams.MinimumIncome
		}
	}

	return minimumIncomeDefault
}

func (ref *minimumIncomeAnalyzer) validate() ([]entity.RuleResultV2, error) {
	if ref.person.PersonType != values.PersonTypeIndividual {
		return nil, errors.New("Minimum income rule is only applied to individual")
	}

	if ref.person.Individual == nil {
		return nil, errors.New("Individual is required for minimum income rule")
	}

	return nil, nil
}

func (ref *minimumIncomeAnalyzer) Name() values.RuleSet {
	return values.RuleSetMinimumIncome
}

package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type highRiskActivity struct {
	person                  entity.Person
	economicActivityService _interfaces.EconomicActivityService
}

func NewHighRiskActivityAnalyzer(person entity.Person,
	economicalActivityService _interfaces.EconomicActivityService) entity.Rule {
	return &highRiskActivity{
		person:                  person,
		economicActivityService: economicalActivityService,
	}
}

func (ref *highRiskActivity) Analyze() ([]entity.RuleResultV2, error) {
	riskyActivity := entity.NewRuleResultV2(values.RuleSetActivityRisk, values.RuleNameHighRiskActivity)

	if ref.person.PersonType != values.PersonTypeCompany {
		return nil, errors.WithStack(errors.Errorf("profile %s is not a company", ref.person.EntityID))
	}

	if ref.person.EnrichedInformation == nil {
		metadata, _ := json.Marshal(fmt.Sprintf("company (%s) not found in bureau", ref.person.DocumentNumber))
		riskyActivity.SetResult(values.ResultStatusRejected).SetMetadata(metadata)
		return []entity.RuleResultV2{*riskyActivity}, nil
	}

	economicActivityCode := ref.person.EnrichedInformation.EconomicActivity
	economicActivity, exists, err := ref.economicActivityService.Get(economicActivityCode)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if exists && economicActivity.RiskValue {
		metadata, _ := json.Marshal(fmt.Sprintf("economic activity (%s) is high risk", economicActivityCode))
		activityCode := ActivityCode{Code: economicActivityCode}
		riskyActivity.SetResult(values.ResultStatusAnalysing).SetPending(true).
			SetMetadata(metadata).AddProblem(values.ProblemCodeEconomicalActivityRiskHigh, activityCode)
		return []entity.RuleResultV2{*riskyActivity}, nil
	}

	riskyActivity.SetResult(values.ResultStatusApproved)
	return []entity.RuleResultV2{*riskyActivity}, nil

}

func (ref *highRiskActivity) Name() values.RuleSet {
	return values.RuleSetActivityRisk
}

type ActivityCode struct {
	Code string `json:"activity_code"`
}

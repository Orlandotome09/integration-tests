package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
)

const (
	statusRegular = "REGULAR"
	statusPending = "PENDING_REGULARIZATION"
)

type bureauAnalyzer struct {
	person                     entity.Person
	approvedStatuses           []string
	notFoundInSerasaStatus     *values.Result
	notFoundInSerasaPending    *bool
	hasProblemsInSerasaStatus  *values.Result
	hasProblemsInSerasaPending *bool
}

func NewBureauAnalyzer(person entity.Person,
	approvedStatuses []string,
	notFoundInSerasaStatus *values.Result,
	notFoundInSerasaPending *bool,
	hasProblemsInSerasaStatus *values.Result,
	hasProblemsInSerasaPending *bool) entity.Rule {
	return &bureauAnalyzer{
		person:                     person,
		approvedStatuses:           approvedStatuses,
		notFoundInSerasaStatus:     notFoundInSerasaStatus,
		notFoundInSerasaPending:    notFoundInSerasaPending,
		hasProblemsInSerasaStatus:  hasProblemsInSerasaStatus,
		hasProblemsInSerasaPending: hasProblemsInSerasaPending,
	}
}

func (ref *bureauAnalyzer) Analyze() ([]entity.RuleResultV2, error) {

	notFoundInSerasaResult, found := ref.analyzeNotFoundRule()
	if !found {
		hasProblemsInSerasaResult := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa)
		return []entity.RuleResultV2{*notFoundInSerasaResult, *hasProblemsInSerasaResult}, nil
	}

	hasProblemsInSerasaResult := ref.analyzeStatusRule()

	return []entity.RuleResultV2{*notFoundInSerasaResult, *hasProblemsInSerasaResult}, nil
}

func (ref *bureauAnalyzer) analyzeNotFoundRule() (ruleResult *entity.RuleResultV2, found bool) {

	notFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa)

	if ref.person.EnrichedInformation != nil {
		notFoundInSerasa.SetResult(values.ResultStatusApproved).SetPending(false)
		return notFoundInSerasa, true
	}

	notFoundInSerasa.SetResult(values.ResultStatusRejected).
		SetPending(false).
		AddProblem(values.ProblemCodeNotFoundAtBureau, map[string]interface{}{"entity_id": ref.person.EntityID, "document_number": ref.person.DocumentNumber})

	if ref.notFoundInSerasaStatus != nil {
		notFoundInSerasa.SetResult(*ref.notFoundInSerasaStatus)
	}
	if ref.notFoundInSerasaPending != nil {
		notFoundInSerasa.SetPending(*ref.notFoundInSerasaPending)
	}

	return notFoundInSerasa, false

}

func (ref *bureauAnalyzer) analyzeStatusRule() *entity.RuleResultV2 {

	hasProblemsInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa)

	if ref.person.EnrichedInformation.BureauStatus == statusRegular ||
		ref.person.EnrichedInformation.BureauStatus == statusPending ||
		ref.isApprovedStatus(ref.person.EnrichedInformation.BureauStatus) {

		hasProblemsInSerasa.SetResult(values.ResultStatusApproved).SetPending(false)
		return hasProblemsInSerasa
	}

	msg := fmt.Sprintf("Bureau status not regular for profile %v. Found BureauStatus: %v",
		ref.person.EntityID, ref.person.EnrichedInformation.BureauStatus)
	metadata, _ := json.Marshal(msg)

	hasProblemsInSerasa.SetResult(values.ResultStatusRejected).
		SetPending(false).SetMetadata(metadata).
		AddProblem(values.ProblemCodeBureauStatusNotRegular, map[string]interface{}{"profile_id": ref.person.EntityID, "status": ref.person.EnrichedInformation.BureauStatus})

	if ref.hasProblemsInSerasaStatus != nil {
		hasProblemsInSerasa.SetResult(*ref.hasProblemsInSerasaStatus)
	}

	if ref.hasProblemsInSerasaPending != nil {
		hasProblemsInSerasa.SetPending(*ref.hasProblemsInSerasaPending)
	}

	return hasProblemsInSerasa

}

func (ref *bureauAnalyzer) Name() values.RuleSet {
	return values.RuleSetSerasa
}

func (ref *bureauAnalyzer) isApprovedStatus(status string) bool {
	if ref.person.PersonType == values.PersonTypeIndividual {
		for _, approvedStatus := range ref.approvedStatuses {
			if status == approvedStatus {
				return true
			}
		}
	}

	return false
}

package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyze(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	person := entity.Person{ProfileID: profileID, DocumentNumber: documentNumber, EnrichedInformation: &entity.EnrichedInformation{BureauStatus: statusRegular}}
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{},
		nil,
		nil,
		nil,
		nil)

	expectedNotFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa).
		SetResult(values.ResultStatusApproved)
	expectedIrregularInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa).
		SetResult(values.ResultStatusApproved)

	rulesResult, err := bureauAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, *expectedNotFoundInSerasa, rulesResult[0])
	assert.Equal(t, *expectedIrregularInSerasa, rulesResult[1])
}

func TestAnalyze_statusNotFound(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	person := entity.Person{ProfileID: profileID, DocumentNumber: documentNumber, EntityID: profileID}
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{},
		nil,
		nil,
		nil,
		nil)

	expectedNotFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa).
		SetResult(values.ResultStatusRejected).SetPending(false).
		AddProblem(values.ProblemCodeNotFoundAtBureau, map[string]interface{}{"entity_id": profileID, "document_number": documentNumber})
	expectedIrregularInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa).
		SetResult(values.ResultStatusIgnored)

	rulesResult, err := bureauAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, *expectedNotFoundInSerasa, rulesResult[0])
	assert.Equal(t, *expectedIrregularInSerasa, rulesResult[1])
}

func TestAnalyze_should_return_the_result_configured_for_status_not_found(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	person := entity.Person{ProfileID: profileID, DocumentNumber: documentNumber, EntityID: profileID}
	notFoundInSerasaStatus := values.ResultStatusApproved
	notFoundInSerasaPending := false
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{},
		&notFoundInSerasaStatus,
		&notFoundInSerasaPending,
		nil,
		nil)

	expectedNotFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa).
		SetResult(notFoundInSerasaStatus).SetPending(notFoundInSerasaPending).
		AddProblem(values.ProblemCodeNotFoundAtBureau, map[string]interface{}{"entity_id": profileID, "document_number": documentNumber})
	expectedIrregularInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa).
		SetResult(values.ResultStatusIgnored)

	rulesResult, err := bureauAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, *expectedNotFoundInSerasa, rulesResult[0])
	assert.Equal(t, *expectedIrregularInSerasa, rulesResult[1])
}

func TestAnalyze_should_return_the_result_configured_for_status_not_regular(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	person := entity.Person{ProfileID: profileID, EntityID: profileID, DocumentNumber: documentNumber, EnrichedInformation: &entity.EnrichedInformation{BureauStatus: "NOT_REGULAR"}}
	hasProblemsInSerasaStatus := values.ResultStatusApproved
	hasProblemsInSerasaPending := false
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{},
		nil,
		nil,
		&hasProblemsInSerasaStatus,
		&hasProblemsInSerasaPending)

	expectedNotFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa).
		SetResult(values.ResultStatusApproved)
	msg := fmt.Sprintf("Bureau status not regular for profile %v. Found BureauStatus: %v",
		profileID, person.EnrichedInformation.BureauStatus)
	metadata, _ := json.Marshal(msg)
	expectedIrregularInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa).
		SetResult(hasProblemsInSerasaStatus).SetPending(hasProblemsInSerasaPending).SetMetadata(metadata).
		AddProblem(values.ProblemCodeBureauStatusNotRegular, map[string]interface{}{"profile_id": profileID, "status": person.EnrichedInformation.BureauStatus})

	rulesResult, err := bureauAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, *expectedNotFoundInSerasa, rulesResult[0])
	assert.Equal(t, *expectedIrregularInSerasa, rulesResult[1])
}

func TestAnalyze_should_return_status_pending_regularization(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	person := entity.Person{ProfileID: profileID, EntityID: profileID, DocumentNumber: documentNumber, EnrichedInformation: &entity.EnrichedInformation{BureauStatus: "PENDING_REGULARIZATION"}}
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{},
		nil,
		nil,
		nil,
		nil)

	expectedNotFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa).
		SetResult(values.ResultStatusApproved)
	expectedIrregularInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa).
		SetResult(values.ResultStatusApproved).SetPending(false)

	rulesResult, err := bureauAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, *expectedNotFoundInSerasa, rulesResult[0])
	assert.Equal(t, *expectedIrregularInSerasa, rulesResult[1])
}

func TestAnalyze_should_return_status_pending_regularization_but_approved_by_config(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	bureauStatusPendingRegularization := "PENDING_REGULARIZATION"
	person := entity.Person{ProfileID: profileID, EntityID: profileID, DocumentNumber: documentNumber, PersonType: values.PersonTypeIndividual, EnrichedInformation: &entity.EnrichedInformation{BureauStatus: bureauStatusPendingRegularization}}
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{bureauStatusPendingRegularization},
		nil,
		nil,
		nil,
		nil)

	expectedNotFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa).
		SetResult(values.ResultStatusApproved)
	expectedIrregularInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa).
		SetResult(values.ResultStatusApproved)

	rulesResult, err := bureauAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, *expectedNotFoundInSerasa, rulesResult[0])
	assert.Equal(t, *expectedIrregularInSerasa, rulesResult[1])
}

func TestAnalyze_statusNotRegular(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	person := entity.Person{ProfileID: profileID, EntityID: profileID, DocumentNumber: documentNumber, EnrichedInformation: &entity.EnrichedInformation{BureauStatus: "NOT_REGULAR"}}
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{},
		nil,
		nil,
		nil,
		nil)

	expectedNotFoundInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameNotFoundInSerasa).
		SetResult(values.ResultStatusApproved)
	msg := fmt.Sprintf("Bureau status not regular for profile %v. Found BureauStatus: %v",
		profileID, person.EnrichedInformation.BureauStatus)
	metadata, _ := json.Marshal(msg)
	expectedIrregularInSerasa := entity.NewRuleResultV2(values.RuleSetSerasa, values.RuleNameHasProblemsInSerasa).
		SetResult(values.ResultStatusRejected).SetPending(false).SetMetadata(metadata).
		AddProblem(values.ProblemCodeBureauStatusNotRegular, map[string]interface{}{"profile_id": profileID, "status": person.EnrichedInformation.BureauStatus})

	rulesResult, err := bureauAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, *expectedNotFoundInSerasa, rulesResult[0])
	assert.Equal(t, *expectedIrregularInSerasa, rulesResult[1])
}

func TestGetBureauName(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "123"
	person := entity.Person{ProfileID: profileID, DocumentNumber: documentNumber}
	bureauAnalyzer := NewBureauAnalyzer(person,
		[]string{},
		nil,
		nil,
		nil,
		nil)

	received := bureauAnalyzer.Name()

	assert.Equal(t, values.RuleSetSerasa, received)
}

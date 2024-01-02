package model

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"testing"
	"time"
)

var modelState State
var domainState entity2.State

func prepare() {

	entityId := uuid.New()
	timeNow, _ := time.Parse(time.ANSIC, "Mon Jan _1 15:04:05 2020")

	metadataString := `["Document File REGISTRATION_FORM Not Found"]`
	tags := []string{"Some tag", "Another tag"}

	metadataBytes, _ := json.Marshal(metadataString)
	tagsBytes, _ := json.Marshal(tags)

	modelRuleResult := RuleResultV2{
		RuleSet:  "IS_UNDER_AGE",
		RuleName: "CUSTOMER_IS_UNDER_AGE",
		Result:   "APPROVED",
		ExpireAt: &timeNow,
		Metadata: metadataBytes,
		Pending:  false,
		Tags:     tagsBytes,
		Problems: nil,
	}
	modelRuleResult2 := RuleResultV2{
		RuleSet:  "INCOMPLETE",
		RuleName: "BLOCKED",
		Result:   "APPROVED",
		ExpireAt: &timeNow,
		Metadata: metadataBytes,
		Pending:  false,
		Tags:     tagsBytes,
		Problems: nil,
	}

	modelValidationStepsResults := []ValidationStepResult{
		{
			StepNumber: 0,
			Result:     "APPROVED",
			RuleResults: []RuleResultV2{
				modelRuleResult,
				modelRuleResult2,
			},
		},
	}

	modelState = State{
		EntityID:               entityId,
		EngineName:             "PROFILE",
		Result:                 "APPROVED",
		ValidationStepsResults: modelValidationStepsResults,
		RuleNames:              []string{"CUSTOMER_IS_UNDER_AGE", "CONSUMER"},
		Pending:                false,
		ExecutionTime:          timeNow,
		CreatedAt:              timeNow,
		UpdatedAt:              timeNow,
	}

	ruleResultsV2 := []entity2.RuleResultV2{
		{
			RuleSet:  values2.RuleSetIsUnderAge,
			RuleName: values2.RuleNameCustomerIsUnderAge,
			Result:   values2.ResultStatusApproved,
			ExpireAt: &timeNow,
			Metadata: metadataBytes,
			Pending:  false,
			Tags:     tags,
			Problems: []entity2.Problem{},
		},
		{
			RuleSet:  values2.RuleSetIncomplete,
			RuleName: values2.RuleNameBlocked,
			Result:   values2.ResultStatusApproved,
			ExpireAt: &timeNow,
			Metadata: metadataBytes,
			Pending:  false,
			Tags:     tags,
			Problems: []entity2.Problem{},
		},
	}

	domainValidationStepResults := []entity2.ValidationStepResult{
		{
			StepNumber:      0,
			Result:          values2.ResultStatusApproved,
			SkipForApproval: false,
			RuleResults:     ruleResultsV2,
		},
	}

	domainState = entity2.State{
		EntityID:               entityId,
		EngineName:             values2.EngineNameProfile,
		Result:                 values2.ResultStatusApproved,
		ExecutionTime:          timeNow,
		CreatedAt:              timeNow,
		UpdatedAt:              timeNow,
		ValidationStepsResults: domainValidationStepResults,
		Pending:                false,
		RuleNames:              []values2.RuleName{values2.RuleName("CUSTOMER_IS_UNDER_AGE"), values2.RuleName("CONSUMER")},
	}

}

func Test_State_Domain_To_Model(t *testing.T) {

	prepare()

	result := modelState.ToDomain()

	assert.Equal(t, result, domainState)

}

func Test_State_Model_To_Domain(t *testing.T) {

	prepare()

	result := StateFromDomain(&domainState)

	assert.Equal(t, result, modelState)

}

func Test_State_Domain_To_Model_and_Back_To_Domain(t *testing.T) {

	prepare()

	modelStateAux := StateFromDomain(&domainState)

	result := modelStateAux.ToDomain()

	assert.Equal(t, result, domainState)

}

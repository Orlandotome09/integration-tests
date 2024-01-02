package preconditionContractRule

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	stateService     *mocks.StateService
	preconditionRule entity.Rule
)

func TestMain(m *testing.M) {
	stateService = &mocks.StateService{}
	os.Exit(m.Run())
}

func TestAnalyze_StateNotFound(t *testing.T) {
	contractID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ContractID: &contractID, ProfileID: &profileID}
	preconditionRule = New(contract, stateService)

	var state *entity.State = nil
	stateService.On("Get", *contract.ProfileID).Return(state, false, nil)

	expectedMetadata, _ := getStateNotFoundMetadata(*contract.ProfileID)
	expectedResult := []entity.RuleResultV2{
		*entity.NewRuleResultV2(values.RuleSetPreconditionContract, values.RuleNameProfileApproved).
			SetResult(values.ResultStatusAnalysing).SetMetadata(expectedMetadata).
			AddProblem(values.ProblemCodeStateNotFound, contract.ProfileID.String()),
	}

	results, err := preconditionRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, results)
}

func TestAnalyze_Approved(t *testing.T) {
	contractID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ContractID: &contractID, ProfileID: &profileID}
	preconditionRule = New(contract, stateService)

	state := &entity.State{Result: ApprovedStatus}
	stateService.On("Get", *contract.ProfileID).Return(state, true, nil)

	expectedResult := []entity.RuleResultV2{
		*entity.NewRuleResultV2(values.RuleSetPreconditionContract, values.RuleNameProfileApproved).
			SetResult(values.ResultStatusApproved),
	}

	results, err := preconditionRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, results)
}

func TestAnalyze_Analysing(t *testing.T) {
	contractID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ContractID: &contractID, ProfileID: &profileID}
	preconditionRule = New(contract, stateService)

	state := &entity.State{Result: values.ResultStatusAnalysing}
	stateService.On("Get", *contract.ProfileID).Return(state, true, nil)

	metadata, _ := json.Marshal(fmt.Sprintf("Profile is not approved: %v", profileID))

	expectedResult := []entity.RuleResultV2{
		*entity.NewRuleResultV2(values.RuleSetPreconditionContract, values.RuleNameProfileApproved).
			SetResult(values.ResultStatusAnalysing).SetMetadata(metadata).
			AddProblem(values.ProblemCodeProfileNotApproved, profileID.String()),
	}

	results, err := preconditionRule.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, results)
}

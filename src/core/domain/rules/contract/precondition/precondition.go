package preconditionContractRule

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/contract"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

var (
	ApprovedStatus = values.ResultStatusApproved
)

type preconditionContractRule struct {
	contractRule.ContractRule
	stateService interfaces.StateService
}

func New(contract entity.Contract, stateService interfaces.StateService) entity.Rule {
	return &preconditionContractRule{
		ContractRule: contractRule.ContractRule{Contract: contract},
		stateService: stateService,
	}
}

func (ref *preconditionContractRule) Analyze() ([]entity.RuleResultV2, error) {
	profileApproved := entity.NewRuleResultV2(values.RuleSetPreconditionContract, values.RuleNameProfileApproved)
	profileApproved.SetResult(values.ResultStatusAnalysing)

	state, exists, err := ref.stateService.Get(*ref.Contract.ProfileID)
	if err != nil {
		return nil, err
	}

	if !exists {
		metadata, _ := getStateNotFoundMetadata(*ref.Contract.ProfileID)
		profileApproved.
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeStateNotFound, ref.Contract.ProfileID.String())
	} else {
		if state.Result == ApprovedStatus {
			profileApproved.SetResult(ApprovedStatus)
		} else {
			metadata, _ := getProfileNotApprovedMetadata(*ref.Contract.ProfileID)
			profileApproved.
				SetMetadata(metadata).
				AddProblem(values.ProblemCodeProfileNotApproved, ref.Contract.ProfileID.String())
		}
	}

	return []entity.RuleResultV2{*profileApproved}, nil
}

func (ref *preconditionContractRule) Name() values.RuleSet {
	return values.RuleSetPreconditionContract
}

func getStateNotFoundMetadata(profileID uuid.UUID) ([]byte, error) {
	return json.Marshal(fmt.Sprintf("State not found for profile: %v",
		profileID))
}

func getProfileNotApprovedMetadata(profileID uuid.UUID) ([]byte, error) {
	return json.Marshal(fmt.Sprintf("Profile is not approved: %v",
		profileID))
}

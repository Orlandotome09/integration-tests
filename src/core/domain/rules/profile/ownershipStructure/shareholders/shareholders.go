package shareholders

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type shareholdersRule struct {
	profile         entity.Profile
	personProcessor interfaces.CompliancePersonProcessor
}

func NewShareholdersRule(profile entity.Profile,
	personProcessor interfaces.CompliancePersonProcessor) interfaces.ShareholdersAnalyzer {
	return &shareholdersRule{
		profile:         profile,
		personProcessor: personProcessor,
	}
}

func (ref *shareholdersRule) Analyze(ownershipStructure entity.OwnershipStructure) (*entity.RuleResultV2, error) {
	shareholdersResult := entity.NewRuleResultV2(values.RuleSetOwnershipStructure, values.RuleNameShareholders).SetResult(values.ResultStatusIgnored)

	result := make([]string, 0)
	notApprovedShareholders := make([]map[string]interface{}, 0)
	hasRejected := false

	for _, shareholder := range ownershipStructure.Shareholders {
		shareholder.Person.PartnerID = ref.profile.Person.PartnerID
		shareholder.Person.ProfileID = ref.profile.Person.ProfileID
		shareholder.Person.OfferType = ref.profile.Person.OfferType

		state, err := ref.personProcessor.ExecuteForPerson(shareholder.Person, ref.profile.Person.OfferType)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if state == nil {
			return nil, errors.New(fmt.Sprintf("[OwnershipStructureRule]State returned is nil for shareholder %+v and profile %+v", shareholder.ShareholderID, ref.profile.ProfileID))
		}

		if state.Result != values.ResultStatusApproved {
			if state.Result == values.ResultStatusRejected {
				hasRejected = true
			}
			result = append(result, fmt.Sprintf("Shareholder with Document Number %v is not Approved", shareholder.Person.DocumentNumber))
			notApprovedShareholders = append(notApprovedShareholders, map[string]interface{}{
				"document_number": shareholder.Person.DocumentNumber,
				"shareholder_id":  shareholder.ShareholderID.String(),
			})
		}
	}

	if len(result) > 0 {
		metadata, err := json.Marshal(result)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		shareholdersResult.
			SetResult(values.ResultStatusAnalysing).
			SetPending(true).
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeShareholderNotApproved, notApprovedShareholders)

		if hasRejected {
			shareholdersResult.SetResult(values.ResultStatusRejected).SetPending(false)
		}

		return shareholdersResult, nil
	}

	shareholdersResult.SetResult(values.ResultStatusApproved)

	return shareholdersResult, nil
}

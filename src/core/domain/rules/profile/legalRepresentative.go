package profile

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type legalRepresentativesRule struct {
	ProfileRule
	personProcessor interfaces.CompliancePersonProcessor
}

func NewLegalRepresentativeRule(profile entity.Profile,
	personProcessor interfaces.CompliancePersonProcessor) entity.Rule {
	return &legalRepresentativesRule{
		ProfileRule: ProfileRule{
			profile: profile,
		},
		personProcessor: personProcessor,
	}
}

func (ref *legalRepresentativesRule) Analyze() ([]entity.RuleResultV2, error) {
	ruleLegalRepresentativesResult := entity.NewRuleResultV2(values.RuleSetLegalRepresentatives, values.RuleNameLegalRepresentativesResult)

	result := make([]string, 0)
	notApprovedIDs := make([]string, 0)
	hasRejected := false

	for _, lr := range ref.profile.LegalRepresentatives {
		lr.PartnerID = ref.profile.PartnerID
		lr.OfferType = ref.profile.Person.OfferType
		state, err := ref.personProcessor.ExecuteForPerson(lr.Person, ref.profile.OfferType)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if state == nil {
			return nil, errors.New(fmt.Sprintf("[LegalRepresentativeRule]State returned is nil for legalRepresentative %+v and profile %+v", lr.LegalRepresentativeID, ref.profile.ProfileID))
		}

		if state.Result != values.ResultStatusApproved {
			if state.Result == values.ResultStatusRejected {
				hasRejected = true
			}
			result = append(result, fmt.Sprintf("Legal Representative %v is not Approved", lr.LegalRepresentativeID))
			notApprovedIDs = append(notApprovedIDs, lr.LegalRepresentativeID.String())
		}
	}

	if len(result) > 0 {
		metadata, err := json.Marshal(result)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		ruleLegalRepresentativesResult.
			SetResult(values.ResultStatusAnalysing).
			SetPending(true).
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeLegalRepresentativeNotApproved, notApprovedIDs)

		if hasRejected {
			ruleLegalRepresentativesResult.SetResult(values.ResultStatusRejected).SetPending(false)
		}

		return []entity.RuleResultV2{*ruleLegalRepresentativesResult}, nil
	}

	ruleLegalRepresentativesResult.SetResult(values.ResultStatusApproved).SetPending(false)
	return []entity.RuleResultV2{*ruleLegalRepresentativesResult}, nil
}

func (ref *legalRepresentativesRule) Name() values.RuleSet {
	return values.RuleSetLegalRepresentatives
}

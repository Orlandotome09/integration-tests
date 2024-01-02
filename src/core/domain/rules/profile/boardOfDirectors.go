package profile

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var validLegalNatures = []string{
	"2143", // Cooperativa
	"2046", // Sociedade Anonima Aberta
	"1210", // Associação Pública
	"3204", // Estabelecimento, no Brasil, de Fundação ou Associação Estrangeiras
	"3212", // Fundação ou Associação domiciliada no exterior
	"3999", // Associação Privada
}

type boardOfDirectorsRule struct {
	ProfileRule
	personProcessor interfaces.CompliancePersonProcessor
}

func NewBoardOfDirectorsRule(profile entity.Profile,
	personProcessor interfaces.CompliancePersonProcessor) entity.Rule {
	return &boardOfDirectorsRule{
		ProfileRule: ProfileRule{
			profile: profile,
		},
		personProcessor: personProcessor,
	}
}

func (ref *boardOfDirectorsRule) Analyze() ([]entity.RuleResultV2, error) {
	ruleBoardOfDirectorsCompleteResult := entity.NewRuleResultV2(values.RuleSetBoardOfDirectors, values.RuleNameBoardOfDirectorsComplete)
	ruleBoardOfDirectorsResult := entity.NewRuleResultV2(values.RuleSetBoardOfDirectors, values.RuleNameBoardOfDirectorsResult)

	if !ref.shouldAnalyzeRule() {
		// Legal Nature is not subject of analysis for this Rule, or it is empty
		ruleBoardOfDirectorsCompleteResult.SetResult(values.ResultStatusIgnored).SetPending(false)
		ruleBoardOfDirectorsResult.SetResult(values.ResultStatusIgnored).SetPending(false)
		return []entity.RuleResultV2{*ruleBoardOfDirectorsCompleteResult, *ruleBoardOfDirectorsResult}, nil
	}

	var boardOfDirectors []entity.Director = nil

	if ref.profile.EnrichedInformation != nil {
		boardOfDirectors = ref.profile.EnrichedInformation.BoardOfDirectors
	}

	if len(boardOfDirectors) == 0 {
		boardOfDirectors = ref.profile.BoardOfDirectors
	}

	if len(boardOfDirectors) == 0 {
		ruleBoardOfDirectorsCompleteResult.
			SetResult(values.ResultStatusAnalysing).
			SetPending(true).
			SetMetadata(getBoardOfDirectorsIncompleteMetadata(*ref.profile.ProfileID)).
			AddProblem(values.ProblemCodeBoardOfDirectorsIncomplete, ref.profile.ProfileID.String())
		ruleBoardOfDirectorsResult.SetResult(values.ResultStatusIgnored).SetPending(false)
		return []entity.RuleResultV2{*ruleBoardOfDirectorsCompleteResult, *ruleBoardOfDirectorsResult}, nil
	}

	ruleBoardOfDirectorsCompleteResult.SetResult(values.ResultStatusApproved).SetPending(false)

	result := make([]string, 0)
	notApprovedIDs := make([]string, 0)
	hasRejected := false
	for _, director := range boardOfDirectors {
		director.Person.PartnerID = ref.profile.PartnerID
		director.Person.OfferType = ref.profile.Person.OfferType
		state, err := ref.personProcessor.ExecuteForPerson(director.Person, ref.profile.OfferType)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if state == nil {
			return nil, errors.New(fmt.Sprintf("[BoardOfDirectorsRule]State returned is nil for director %+v and profile %+v", director.DirectorID, ref.profile.ProfileID))
		}

		if state.Result != values.ResultStatusApproved {
			if state.Result == values.ResultStatusRejected {
				hasRejected = true
			}
			result = append(result, fmt.Sprintf("Director %v is not Approved", director.DirectorID))
			notApprovedIDs = append(notApprovedIDs, director.DirectorID.String())
		}
	}

	if len(result) > 0 {
		metadata, err := json.Marshal(result)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ruleBoardOfDirectorsResult.
			SetResult(values.ResultStatusAnalysing).
			SetPending(true).
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeDirectorNotApproved, notApprovedIDs)

		if hasRejected {
			ruleBoardOfDirectorsResult.SetResult(values.ResultStatusRejected).SetPending(false)
		}

		return []entity.RuleResultV2{*ruleBoardOfDirectorsCompleteResult, *ruleBoardOfDirectorsResult}, nil
	}

	ruleBoardOfDirectorsResult.SetResult(values.ResultStatusApproved).SetPending(false)
	return []entity.RuleResultV2{*ruleBoardOfDirectorsCompleteResult, *ruleBoardOfDirectorsResult}, nil
}

func (ref *boardOfDirectorsRule) Name() values.RuleSet {
	return values.RuleSetBoardOfDirectors
}

func (ref *boardOfDirectorsRule) shouldAnalyzeRule() bool {
	legalNature := ""

	// Use enriched LegalNature if exists
	if ref.profile.EnrichedInformation != nil {
		legalNature = ref.profile.EnrichedInformation.LegalNature
	}

	// Use profile LegalNature if there is no enriched info
	if legalNature == "" && ref.profile.Company != nil {
		legalNature = ref.profile.Company.LegalNature
	}

	found := false
	for _, validLegalNature := range validLegalNatures {
		if validLegalNature == legalNature {
			found = true
		}
	}

	return found
}

func getBoardOfDirectorsIncompleteMetadata(profileID uuid.UUID) []byte {
	message, _ := json.Marshal(fmt.Sprintf("Board of Directors not found for profile: %v",
		profileID))
	return message
}

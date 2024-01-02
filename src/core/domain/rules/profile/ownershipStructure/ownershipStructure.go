package ownershipStructure

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

var invalidLegalNatures = []string{
	"2143", // Cooperativa
	"2046", // Sociedade Anonima Aberta
	"1210", // Associação Pública
	"3204", // Estabelecimento, no Brasil, de Fundação ou Associação Estrangeiras
	"3212", // Fundação ou Associação domiciliada no exterior
	"3999", // Associação Privada
}

type ownershipStructureRule struct {
	profile          entity.Profile
	shareholdingRule interfaces.ShareholdingAnalyzer
	shareholdersRule interfaces.ShareholdersAnalyzer
}

func NewOwnershipStructureRule(profile entity.Profile,
	shareholdingRule interfaces.ShareholdingAnalyzer,
	shareholdersRule interfaces.ShareholdersAnalyzer) entity.Rule {
	return &ownershipStructureRule{
		profile:          profile,
		shareholdingRule: shareholdingRule,
		shareholdersRule: shareholdersRule,
	}
}

func (ref *ownershipStructureRule) Analyze() ([]entity.RuleResultV2, error) {
	shareholdersResult := entity.NewRuleResultV2(values.RuleSetOwnershipStructure, values.RuleNameShareholders).SetResult(values.ResultStatusIgnored)
	shareholdingResult := entity.NewRuleResultV2(values.RuleSetOwnershipStructure, values.RuleNameShareholding).SetResult(values.ResultStatusIgnored)

	if !ref.shouldAnalyzeRule() {
		return []entity.RuleResultV2{*shareholdingResult, *shareholdersResult}, nil
	}

	shareholdingResult, ownershipStructure, err := ref.shareholdingRule.Analyze()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if ownershipStructure == nil {
		return []entity.RuleResultV2{*shareholdingResult, *shareholdersResult}, nil
	}

	if shareholdingResult.Result != values.ResultStatusApproved {
		existsApprovedOverride := ref.hasApprovedOverrideForShareholdingRule()

		if !existsApprovedOverride {
			return []entity.RuleResultV2{*shareholdingResult, *shareholdersResult}, nil
		}
	}

	shareholdersResult, err = ref.shareholdersRule.Analyze(*ownershipStructure)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return []entity.RuleResultV2{*shareholdingResult, *shareholdersResult}, nil
}

func (ref *ownershipStructureRule) Name() values.RuleSet {
	return values.RuleSetOwnershipStructure
}

func (ref *ownershipStructureRule) hasApprovedOverrideForShareholdingRule() bool {

	for _, override := range ref.profile.Overrides {
		if override.RuleSet == values.RuleSetOwnershipStructure &&
			override.RuleName == values.RuleNameShareholding &&
			override.Result == values.ResultStatusApproved {
			return true
		}
	}

	return false
}

func (ref *ownershipStructureRule) shouldAnalyzeRule() bool {
	legalNature := ""

	// Use enriched LegalNature if exists
	if ref.profile.EnrichedInformation != nil {
		legalNature = ref.profile.EnrichedInformation.LegalNature
	}

	// Use profile LegalNature if there is no enriched info
	if legalNature == "" && ref.profile.Company != nil {
		legalNature = ref.profile.Company.LegalNature
	}

	found := true
	for _, invalidLegalNature := range invalidLegalNatures {
		if invalidLegalNature == legalNature {
			found = false
		}
	}

	return found
}

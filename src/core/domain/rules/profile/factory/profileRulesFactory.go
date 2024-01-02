package profileRulesFactory

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/profile"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/profile/ownershipStructure"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/profile/ownershipStructure/shareholders"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/profile/ownershipStructure/shareholding"
	"github.com/sirupsen/logrus"
)

type profileRulesFactory struct {
	personProcessor             interfaces.CompliancePersonProcessor
	personsRulesFactoryInstance interfaces.PersonRulesFactory
}

func New(personProcessor interfaces.CompliancePersonProcessor,
	personsRulesFactoryInstance interfaces.PersonRulesFactory,
) interfaces.ProfileRulesFactory {
	return &profileRulesFactory{
		personProcessor:             personProcessor,
		personsRulesFactoryInstance: personsRulesFactoryInstance,
	}
}

func (ref *profileRulesFactory) GetRules(ruleSetConfig *entity.RuleSetConfig, record *entity.Profile) []entity.Rule {
	var resultRules []entity.Rule

	if ruleSetConfig == nil {
		return nil
	}

	resultRules = ref.personsRulesFactoryInstance.GetRules(*ruleSetConfig, record.Person)

	if ruleSetConfig.OwnershipStructureParams != nil {
		shareholdingRule := shareholding.NewShareholdingRule(*record)
		shareholdersRule := shareholders.NewShareholdersRule(*record, ref.personProcessor)
		resultRules = append(resultRules,
			ownershipStructure.NewOwnershipStructureRule(*record, shareholdingRule, shareholdersRule))
	}

	if ruleSetConfig.LegalRepresentativeParams != nil {
		resultRules = append(resultRules, profile.NewLegalRepresentativeRule(*record, ref.personProcessor))
	}

	if ruleSetConfig.BoardOfDirectorsParams != nil {
		resultRules = append(resultRules, profile.NewBoardOfDirectorsRule(*record, ref.personProcessor))
	}

	if ruleSetConfig.MinimumBillingParams != nil {
		resultRules = append(resultRules, profile.NewMinimumBillingAnalyzer(*record))
	}

	if ruleSetConfig.ManualValidationParams != nil {
		resultRules = append(resultRules, profile.NewManualValidationRule(*record))
	}

	logrus.WithField("ruleSetConfig", ruleSetConfig).
		WithField("rules", resultRules).
		Infof("[profileRulesFactory] Rules created for profile %v", record.ProfileID)

	return resultRules
}

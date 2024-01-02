package _init

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/engines"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/engines/ruleValidator"
	profilerulesfactory "bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/profile/factory"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/statemachine"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/complianceAnalyzer"
	statemanager "bitbucket.org/bexstech/temis-compliance/src/core/useCases/complianceAnalyzer/stateManager"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/personProcessor"
)

func buildProfileRulesFactory() interfaces.ProfileRulesFactory {
	ruleValidatorInstance := ruleValidator.New(statemachine.NewStateMachine())

	personSubEngine := engines.NewPersonSubEngine(ruleValidatorInstance, buildPersonFactory(), personRepository)
	stateManagerSync := statemanager.NewStateManager(buildStateService())
	processorSync := complianceAnalyzer.NewComplianceAnalyzer(overrideRepository, stateManagerSync)
	personProcessor := personProcessor.NewCompliancePersonProcessor(personSubEngine, processorSync)

	return profilerulesfactory.New(personProcessor, buildPersonRulesFactory())

}

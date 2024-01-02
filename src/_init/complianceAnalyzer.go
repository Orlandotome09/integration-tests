package _init

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/complianceAnalyzer"
	statemanager "bitbucket.org/bexstech/temis-compliance/src/core/useCases/complianceAnalyzer/stateManager"
)

func buildComplianceAnalyzer() interfaces.ComplianceAnalyzer {
	stateManager := statemanager.NewStateManager(buildStateService())

	return complianceAnalyzer.NewComplianceAnalyzer(overrideRepository, stateManager)

}

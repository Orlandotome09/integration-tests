package _init

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract"
	contractClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/engines"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/engines/ruleValidator"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/statemachine"
)

func buildContractEngine() interfaces.Engine {

	contractClientInstance := contractClient.New(webClient, temisRegistrationHost)

	contractService := contract.NewContractAdapter(contractClientInstance)

	ruleValidatorInstance := ruleValidator.New(statemachine.NewStateMachine())

	return engines.NewContractEngine(ruleValidatorInstance, contractService, contractRepository, buildContractRulesFactory())
}

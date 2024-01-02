package _init

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile"
	profileHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/engines"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/engines/ruleValidator"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/statemachine"
)

func buildProfileEngine() interfaces.Engine {
	profileClient := profileHttpClient.NewProfileHttpClient(webClient, temisRegistrationHost)
	profileAdapter := profile.NewProfileService(profileClient)

	ruleValidatorInstance := ruleValidator.New(statemachine.NewStateMachine())

	return engines.NewProfileEngine(ruleValidatorInstance,
		buildPosProcessor(),
		profileAdapter,
		buildProfileFactory(),
		profileRepository)

}

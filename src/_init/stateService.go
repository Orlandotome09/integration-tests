package _init

import (
	"context"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract"
	contractClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile"
	profileHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventPublisher/complianceCommandPublisher"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventPublisher/stateEventPublisher"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/state"
	"bitbucket.org/bexstech/temis-compliance/src/infra/queues/pubsubPublisher"
)

func buildStateService() interfaces.StateService {
	ctx := context.Background()

	stateEventsTopicPublisher := pubsubPublisher.New(ctx, bexsDigitalPubsubClient, stateEventsTopic)
	complianceCommandTopicPublisher := pubsubPublisher.New(ctx, bexsPubsubClient, complianceCommandTopic)

	stateEventsPublisher := stateEventPublisher.NewStateEventsPublisher(
		profileRepository, contractRepository, personRepository, stateEventsTopicPublisher, randomIDGenerator)

	complianceCommandsPublisher := complianceCommandPublisher.New(complianceCommandTopicPublisher)

	profileClient := profileHttpClient.NewProfileHttpClient(webClient, temisRegistrationHost)
	contractClientInstance := contractClient.New(webClient, temisRegistrationHost)

	profileAdapter := profile.NewProfileService(profileClient)
	contractService := contract.NewContractAdapter(contractClientInstance)

	return state.NewStateService(stateRepository, profileAdapter, contractService, stateEventsPublisher, complianceCommandsPublisher)
}

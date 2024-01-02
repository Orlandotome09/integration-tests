package _init

import (
	"context"
	"sync"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/complianceCommandProcessor"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichedPersonEventProcessor"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registrationEventProcessor"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventProcessor"
	enginefactory "bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventProcessor/engineFactory"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"
	eventListener "bitbucket.org/bexstech/temis-compliance/src/infra/queues/subscriber"
	"bitbucket.org/bexstech/temis-compliance/src/infra/queues/subscriber/pubsubSubscriber"
	"bitbucket.org/bexstech/temis-compliance/src/infra/queues/subscriber/pubsubSubscriber/mutex"
)

func StartComplianceAsync(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	metricsGenerator := metrics.New()

	mutexInstance := mutex.New()
	complianceCommandSubscriber := pubsubSubscriber.New(&ctx, bexsPubsubClient, complianceCommanndSubscription, metricsGenerator, mutexInstance)
	registrationEventsSubscriber := pubsubSubscriber.New(&ctx, bexsDigitalPubsubClient, registrationEventsSubscription, metricsGenerator, mutexInstance)
	enrichmentEventsSubscriber := pubsubSubscriber.New(&ctx, bexsPubsubClient, complianceSubscriptionToEnrichmentTopic, metricsGenerator, mutexInstance)

	// CnC--------------------------
	engineFactory := enginefactory.NewEngineFactory(buildProfileEngine(), buildContractEngine())

	evtProcessor := eventProcessor.New(buildComplianceAnalyzer(), engineFactory)

	registrationProcessor := registrationEventProcessor.New(evtProcessor)
	enrichedPersonProcessor := enrichedPersonEventProcessor.New(evtProcessor, timeNowGenerator)
	complianceCommandProcessor := complianceCommandProcessor.New(evtProcessor)

	eventListenerInstance := eventListener.New()
	eventListenerInstance.Register(registrationEventsSubscriber, registrationProcessor)
	eventListenerInstance.Register(enrichmentEventsSubscriber, enrichedPersonProcessor)
	eventListenerInstance.Register(complianceCommandSubscriber, complianceCommandProcessor)

	eventListenerInstance.Listen()
}

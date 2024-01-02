package pubsubSubscriber

import (
	"context"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"
	"bitbucket.org/bexstech/temis-compliance/src/infra/queues/subscriber/pubsubSubscriber/message"
	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type queueSubscriber struct {
	ctx            *context.Context
	client         *pubsub.Client
	subscriptionID string
	metrics        interfaces.MetricsGenerator
	mutex          interfaces.Mutex
}

func New(ctx *context.Context, client *pubsub.Client, subscriptionID string, metricsGenerator interfaces.MetricsGenerator, mutex interfaces.Mutex) interfaces.QueueSubscriber {
	return &queueSubscriber{
		ctx:            ctx,
		client:         client,
		subscriptionID: subscriptionID,
		metrics:        metricsGenerator,
		mutex:          mutex,
	}
}

func (ref *queueSubscriber) Subscribe(processor interfaces.MessageProcessor) error {
	sub := ref.client.Subscription(ref.subscriptionID)
	sub.ReceiveSettings.NumGoroutines = 10
	sub.ReceiveSettings.MaxOutstandingMessages = 100

	logrus.Infof("[QueueSubscriber] About to receive messages from Subscription %v", sub.String())

	err := sub.Receive(*ref.ctx, func(_ context.Context, m *pubsub.Message) {
		ref.processMessage(processor, message.New(m, ref.metrics))
	})
	if err != nil {
		logrus.Errorf("[QueueSubscriber] Error receiving messages from Subscription %v. Error: %s", sub.String(), err.Error())

		return err
	}

	return nil
}

func (ref *queueSubscriber) processMessage(processor interfaces.MessageProcessor, message message.Message) {
	result := metrics.EventProcessStatusError

	defer func() {
		message.Nack(result)
	}()

	logrus.
		WithField("message_id", message.ID()).
		WithField("message_data", string(message.Data())).
		Infof("[QueueSubscriber] Processing Compliace Event from Subscription %s", ref.subscriptionID)

	if message.ConcurrencyKey() != "" {
		if !ref.mutex.Lock(message.ConcurrencyKey()) {
			result = metrics.EventProcessStatusPending
			logrus.
				WithField("message_id", message.ID()).
				WithField("message_data", string(message.Data())).
				Infof("[QueueSubscriber] Profile is currently being processed by another routine. Returning to queue")
			return
		}

		defer ref.mutex.Release(message.ConcurrencyKey())
	}

	result, err := processor.Process(message)
	if err != nil {
		logrus.
			WithField("message_id", message.ID()).
			WithField("message_data", string(message.Data())).
			Errorf("Error on processing message: %+v", errors.WithStack(err))
		return
	}

	message.Ack()
}

func (ref *queueSubscriber) SubscriptionName() string {
	return ref.subscriptionID
}

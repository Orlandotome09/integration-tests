package listener

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type pubSubListener struct {
	client         *pubsub.Client
	context        context.Context
	subscriptionID string
}

func NewPubSubListener(client *pubsub.Client, context context.Context, subscriptionID string) Listener {
	return &pubSubListener{
		client:         client,
		context:        context,
		subscriptionID: subscriptionID,
	}
}

func (ref *pubSubListener) processMessage(processor Processor, message *pubsub.Message) {
	defer message.Nack()
	event := values.Event{}
	data := message.Data
	if err := json.Unmarshal(data, &event); err != nil {
		logrus.
			WithField("message_id", message.ID).
			WithField("message_data", string(message.Data)).
			Errorf("Error on unmarshal message: %+v", errors.WithStack(err))
		return
	}

	logrus.Infof("Processing Event: %+v", event)
	if err := processor(&event); err != nil {
		logrus.
			WithField("message_id", message.ID).
			WithField("message_data", string(message.Data)).
			Errorf("Error on processing message: %+v", errors.WithStack(err))
		return
	}
	message.Ack()
}

func (ref *pubSubListener) Listen(processor Processor) error {
	sub := ref.client.Subscription(ref.subscriptionID)
	sub.ReceiveSettings.MaxOutstandingMessages = 50
	sub.ReceiveSettings.MaxExtension = -1

	err := sub.Receive(ref.context, func(_ context.Context, message *pubsub.Message) {
		ref.processMessage(processor, message)
	})
	if err != context.Canceled {
		return errors.WithStack(err)
	}
	return nil
}

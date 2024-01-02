package pubsubPublisher

import (
	"context"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type queuePublisher struct {
	ctx     context.Context
	client  *pubsub.Client
	topicID string
}

func New(ctx context.Context, client *pubsub.Client, topicID string) interfaces.QueuePublisher {
	return &queuePublisher{
		ctx:     ctx,
		client:  client,
		topicID: topicID,
	}
}

func (ref *queuePublisher) Publish(message []byte, orderingKey string, concurrencyKey string) error {
	topic := ref.client.Topic(ref.topicID)

	result := topic.Publish(ref.ctx, &pubsub.Message{
		Data:        message,
		OrderingKey: orderingKey,
		Attributes: map[string]string{
			"concurrency_key": concurrencyKey,
		},
	})

	id, err := result.Get(ref.ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	logrus.
		WithField("message_id", id).
		Infof("[QueuePublisher] Published Message")
	return nil
}

package main

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

const (
	bexsProjectID = "local-project"
)

const (
	stateEventsTopic              = "temis-compliance-state-events-v1"
	stateEventsTopicTestSub       = "temis-compliance-state-events-v1-test-subscription"
	stateEventsTopicTestSubDlq    = "temis-compliance-state-events-v1-test-subscription-dlq"
	stateEventsTopicTestSubDlqSub = "temis-compliance-state-events-v1-test-subscription-dlq-subscription"
)

type pubsubInfo struct {
	ProjectID       string
	TopicID         string
	SubscriptionID  string
	DeadLetterTopic string
}

func main() {
	pubsubs := []pubsubInfo{
		{
			ProjectID:      bexsProjectID,
			TopicID:        "deadletter_temis_event",
			SubscriptionID: "temis_event_to_cnc_deadletter",
		},
		{
			ProjectID:       bexsProjectID,
			TopicID:         "temis-event",
			SubscriptionID:  "temis-event-to-cnc",
			DeadLetterTopic: "projects/local-project/topics/deadletter_temis_event",
		},

		{
			ProjectID:      bexsProjectID,
			TopicID:        "deadletter_compliance_topic",
			SubscriptionID: "notification_subscription_deadletter",
		},
		{
			ProjectID:       bexsProjectID,
			TopicID:         "compliance_topic",
			SubscriptionID:  "notification_subscription",
			DeadLetterTopic: "projects/local-project/topics/deadletter_compliance_topic",
		},
		{
			ProjectID:      bexsProjectID,
			TopicID:        "deadletter_enrichment_topic",
			SubscriptionID: "enrichment_subscription_deadletter",
		},
		{
			ProjectID:       bexsProjectID,
			TopicID:         "temis-enrichment-person-events",
			SubscriptionID:  "temis-enrichment-person-events-subscription",
			DeadLetterTopic: "projects/local-project/topics/deadletter_enrichment_topic",
		},
		{
			ProjectID:      bexsProjectID,
			TopicID:        "tree_adapter_topic",
			SubscriptionID: "tree_adapter_subscription",
		},
		{
			ProjectID:      bexsProjectID,
			TopicID:        "temis-limit",
			SubscriptionID: "limit_subscription",
		},
		{
			ProjectID:      bexsProjectID,
			TopicID:        stateEventsTopicTestSubDlq,
			SubscriptionID: stateEventsTopicTestSubDlqSub,
		},
		{
			ProjectID:       bexsProjectID,
			TopicID:         stateEventsTopic,
			SubscriptionID:  stateEventsTopicTestSub,
			DeadLetterTopic: "projects/local-project/topics/" + stateEventsTopicTestSubDlq,
		},
		{
			ProjectID:      bexsProjectID,
			TopicID:        "temis-registration-events-compliance-subscription-dlq",
			SubscriptionID: "temis-registration-events-compliance-subscription-dlq-subscription",
		},
		{
			ProjectID:       bexsProjectID,
			TopicID:         "temis-registration-events",
			SubscriptionID:  "temis-registration-events-compliance-subscription",
			DeadLetterTopic: "projects/local-project/topics/temis-registration-events-compliance-subscription-dlq",
		},
	}

	for _, i := range pubsubs {
		ctx := context.Background()

		options := []option.ClientOption{}

		client, err := pubsub.NewClient(ctx, i.ProjectID, options...)
		if err != nil {
			logrus.Fatal(err)
		}

		topic := createTopic(i.TopicID, client, ctx)
		createSubscription(i.SubscriptionID, topic, client, ctx, i.DeadLetterTopic)

	}
}

func createTopic(id string, client *pubsub.Client, ctx context.Context) *pubsub.Topic {
	topic := client.Topic(id)
	_ = topic.Delete(ctx)

	topic, err := client.CreateTopic(ctx, id)
	if err != nil {
		logrus.Fatal(err)
	}

	topic.EnableMessageOrdering = true

	logrus.Printf("Topic created: %v\n", topic)

	return topic
}

func createSubscription(id string, topic *pubsub.Topic, client *pubsub.Client, ctx context.Context, deadLetterTopic string) {
	config := pubsub.SubscriptionConfig{
		Topic:                 topic,
		EnableMessageOrdering: false,
		RetainAckedMessages:   false,
		RetryPolicy: &pubsub.RetryPolicy{
			MinimumBackoff: time.Second,
			MaximumBackoff: time.Second * 60,
		},
	}

	if deadLetterTopic != "" {
		config.DeadLetterPolicy = &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     deadLetterTopic,
			MaxDeliveryAttempts: 5,
		}
	}

	s, err := client.CreateSubscription(ctx, id, config)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Printf("Subscription created: %v\n", s)
}

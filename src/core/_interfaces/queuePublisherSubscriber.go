package _interfaces

import "time"

type QueuePublisher interface {
	Publish(message []byte, orderingKey string, concurrencyKey string) error
}

type QueueSubscriber interface {
	Subscribe(processor MessageProcessor) error
	SubscriptionName() string
}

type Message interface {
	ID() string
	Data() []byte
	PublishTime() time.Time
}

package eventListener

import (
	"log"
	"sync"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"github.com/sirupsen/logrus"
)

type eventListener struct {
	eventSubscribers []eventSubscriber
}

type eventSubscriber struct {
	subscriber _interfaces.QueueSubscriber
	processor  _interfaces.MessageProcessor
}

func New() _interfaces.EventListener {
	return &eventListener{
		eventSubscribers: make([]eventSubscriber, 0),
	}
}

func (ref *eventListener) Register(subscriber _interfaces.QueueSubscriber, processor _interfaces.MessageProcessor) {
	ref.eventSubscribers = append(ref.eventSubscribers, eventSubscriber{
		subscriber: subscriber,
		processor:  processor,
	})
}

func (ref *eventListener) Listen() {
	wg := new(sync.WaitGroup)

	for _, e := range ref.eventSubscribers {
		wg.Add(1)
		go func(e eventSubscriber) {
			logrus.Infof("[Subscriber] Started listening to subscription: %s", e.subscriber.SubscriptionName())
			err := e.subscriber.Subscribe(e.processor)
			if err != nil {
				logrus.WithField("error", err.Error()).Errorf("[Subscriber] Cannot read subscription: %s", e.subscriber.SubscriptionName())
				log.Fatal(err.Error())
			}

			logrus.Infof("[Subscriber] Finished listening to subscription: %s", e.subscriber.SubscriptionName())
			wg.Done()
		}(e)
	}
	wg.Wait()
}

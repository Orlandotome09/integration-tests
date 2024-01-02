package listener

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"time"
)

type mockedListener struct {
	TickTime time.Duration
	Event    *values.Event
}

func NewMockedListener(tickTime time.Duration, event *values.Event) Listener {
	return &mockedListener{
		TickTime: tickTime,
		Event:    event,
	}
}

func (ref *mockedListener) Listen(processor Processor) error {
	for {
		time.Sleep(ref.TickTime)
		processor(ref.Event)
	}
}

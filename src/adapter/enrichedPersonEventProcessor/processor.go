package enrichedPersonEventProcessor

import (
	"encoding/json"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"

	"github.com/sirupsen/logrus"
)

type enrichedPersonEventsProcessor struct {
	eventProcessor _interfaces.EventProcessor
	timeGenerator  func() time.Time
}

func New(eventProcessor _interfaces.EventProcessor, timeGenerator func() time.Time) _interfaces.MessageProcessor {
	return &enrichedPersonEventsProcessor{
		eventProcessor: eventProcessor,
		timeGenerator:  timeGenerator,
	}
}

func (ref *enrichedPersonEventsProcessor) Process(message _interfaces.Message) (string, error) {
	enrichedPersonEvent := EnrichedPersonEvent{}
	if err := json.Unmarshal(message.Data(), &enrichedPersonEvent); err != nil {
		logrus.WithField("message_data", string(message.Data())).
			WithField("error", err.Error()).
			Warning("[enrichedPerson EventsProcessor] Error parsing registration event. Discarding broken contract event")

		return metrics.EventProcessStatusDiscarded, nil
	}

	valid, err := enrichedPersonEvent.IsValid()
	if !valid {
		logrus.WithField("message_data", string(message.Data())).
			WithField("error", err.Error()).
			Warning("[enrichedPerson EventsProcessor] Invalid registration event. Discarding event")

		return metrics.EventProcessStatusDiscarded, nil
	}

	event := &values.Event{
		EngineName:  values.EngineNameProfile,
		EventType:   values.EventTypePersonEnriched,
		ParentID:    enrichedPersonEvent.EntityID,
		EntityID:    enrichedPersonEvent.EntityID,
		Date:        ref.timeGenerator(),
		RequestDate: message.PublishTime(),
	}

	_, err = ref.eventProcessor.ExecuteForEvent(event)
	if err != nil {
		return metrics.EventProcessStatusError, err
	}

	return metrics.EventProcessStatusSuccess, nil
}

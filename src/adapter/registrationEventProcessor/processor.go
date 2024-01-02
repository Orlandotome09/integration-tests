package registrationEventProcessor

import (
	"encoding/json"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"

	"github.com/sirupsen/logrus"
)

type registrationEventsProcessor struct {
	eventProcessor _interfaces.EventProcessor
}

func New(eventProcessor _interfaces.EventProcessor) _interfaces.MessageProcessor {
	return &registrationEventsProcessor{
		eventProcessor: eventProcessor,
	}
}

func (ref *registrationEventsProcessor) Process(message _interfaces.Message) (string, error) {
	registrationEvent := RegistrationEvent{}
	if err := json.Unmarshal(message.Data(), &registrationEvent); err != nil {
		logrus.WithField("message_data", string(message.Data())).
			WithField("error", err.Error()).
			Warning("[RegistrationEventsProcessor] Error parsing registration event. Discarding broken contract event")

		return metrics.EventProcessStatusDiscarded, nil
	}

	valid, err := registrationEvent.IsValid()
	if !valid {
		logrus.WithField("message_data", string(message.Data())).
			WithField("error", err.Error()).
			Warning("[RegistrationEventsProcessor] Invalid registration event. Discarding event")

		return metrics.EventProcessStatusDiscarded, nil
	}

	parentID := registrationEvent.ProfileID
	engineName := values.EngineNameProfile

	if registrationEvent.EntityType == values.EntityTypeContract.ToString() {
		parentID = registrationEvent.EntityID
		engineName = values.EngineNameContract
	}

	event := &values.Event{
		EngineName:  engineName,
		ParentID:    parentID,
		Date:        registrationEvent.UpdateDate,
		EventType:   registrationEvent.EventType,
		EntityID:    registrationEvent.EntityID,
		Content:     registrationEvent.Content,
		RequestDate: message.PublishTime(),
	}

	_, err = ref.eventProcessor.ExecuteForEvent(event)
	if err != nil {
		return metrics.EventProcessStatusError, err
	}

	return metrics.EventProcessStatusSuccess, nil
}

package complianceCommandProcessor

import (
	"encoding/json"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"

	"github.com/sirupsen/logrus"
)

type complianceCommandProcessor struct {
	eventProcessor _interfaces.EventProcessor
}

func New(eventProcessor _interfaces.EventProcessor) _interfaces.MessageProcessor {
	return &complianceCommandProcessor{
		eventProcessor: eventProcessor,
	}
}

func (ref *complianceCommandProcessor) Process(message _interfaces.Message) (string, error) {
	command := complianceCommand{}
	if err := json.Unmarshal(message.Data(), &command); err != nil {
		logrus.WithField("message_data", string(message.Data())).
			WithField("error", err.Error()).
			Warning("[ComplianceCommandProcessor] Error parsing command. Discarding broken contract event")

		return metrics.EventProcessStatusDiscarded, nil
	}

	valid, err := command.IsValid()
	if !valid {
		logrus.WithField("message_data", string(message.Data())).
			WithField("error", err.Error()).
			Warning("[ComplianceCommandProcessor] Invalid command. Discarding event")

		return metrics.EventProcessStatusDiscarded, nil
	}

	event := &values.Event{
		EngineName:  command.EngineName,
		EventType:   command.EventType,
		ParentID:    command.EntityID,
		EntityID:    command.EntityID,
		Date:        command.Date,
		Content:     command.Content,
		RequestDate: message.PublishTime(),
	}

	_, err = ref.eventProcessor.ExecuteForEvent(event)
	if err != nil {
		return metrics.EventProcessStatusError, err
	}

	return metrics.EventProcessStatusSuccess, nil
}

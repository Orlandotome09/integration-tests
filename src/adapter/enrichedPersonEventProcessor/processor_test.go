package enrichedPersonEventProcessor

import (
	"encoding/json"
	"testing"
	"time"

	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var now = time.Now()

func timeGenerator() time.Time {
	return now
}

func TestProcess(t *testing.T) {
	eventProcessor := &mocks.EventProcessor{}

	processor := New(eventProcessor, timeGenerator)

	enrichedPersonEvent := &EnrichedPersonEvent{
		EventType:  "PERSON_ENRICHED",
		EntityID:   uuid.New(),
		EntityType: "PROFILE",
		Data: PersonEnrichedData{
			PersonID: uuid.New(),
		},
	}
	b, _ := json.Marshal(enrichedPersonEvent)

	message := &mocks.Message{}
	message.On("Data").Return(b)
	message.On("PublishTime").Return(now)

	event := &values.Event{
		EngineName:  values.EngineNameProfile,
		EventType:   enrichedPersonEvent.EventType,
		ParentID:    enrichedPersonEvent.EntityID,
		EntityID:    enrichedPersonEvent.EntityID,
		Date:        now,
		RequestDate: now,
	}

	eventProcessor.On("ExecuteForEvent", event).Return(nil, nil)

	status, err := processor.Process(message)

	assert.NoError(t, err)
	assert.Equal(t, "success", status)
}

func TestProcess_should_not_process_when_broken_contract(t *testing.T) {
	eventProcessor := &mocks.EventProcessor{}
	registrationEventProcessor := New(eventProcessor, timeGenerator)

	registrationEvent := "not correct contract"

	data, _ := json.Marshal(registrationEvent)
	msg := &mocks.Message{}
	msg.On("Data").Return(data)
	msg.On("PublishTime").Return(now)

	received, err := registrationEventProcessor.Process(msg)

	assert.Nil(t, err)
	assert.Equal(t, metrics.EventProcessStatusDiscarded, received)
	eventProcessor.AssertNumberOfCalls(t, "Process", 0)
}

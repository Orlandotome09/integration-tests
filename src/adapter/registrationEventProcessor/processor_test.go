package registrationEventProcessor

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

func TestProcess(t *testing.T) {
	eventProcessor := &mocks.EventProcessor{}

	processor := New(eventProcessor)

	parentID := uuid.New()

	b, _ := json.Marshal(&RegistrationEvent{
		ProfileID:  uuid.New(),
		ParentID:   &parentID,
		EventType:  "PROFILE_CREATED",
		EntityID:   uuid.New(),
		EntityType: "PROFILE",
		ParentType: "PROFILE",
		UpdateDate: time.Now(),
		Content:    nil,
	})

	now := time.Now()

	message := &mocks.Message{}
	message.On("Data").Return(b)
	message.On("PublishTime").Return(now)

	registrationEvent := RegistrationEvent{}
	json.Unmarshal(b, &registrationEvent)

	event := &values.Event{
		EngineName:  values.EngineNameProfile,
		ParentID:    registrationEvent.ProfileID,
		Date:        registrationEvent.UpdateDate,
		EventType:   registrationEvent.EventType,
		EntityID:    registrationEvent.EntityID,
		Content:     registrationEvent.Content,
		RequestDate: now,
	}

	eventProcessor.On("ExecuteForEvent", event).Return(nil, nil)

	status, err := processor.Process(message)

	assert.NoError(t, err)
	assert.Equal(t, "success", status)
}

func TestProcess_should_not_process_when_broken_contract(t *testing.T) {
	eventProcessor := &mocks.EventProcessor{}
	registrationEventProcessor := New(eventProcessor)

	registrationEvent := "not correct contract"

	data, _ := json.Marshal(registrationEvent)
	msg := &mocks.Message{}
	msg.On("Data").Return(data)
	msg.On("PublishTime").Return(time.Now())

	received, err := registrationEventProcessor.Process(msg)

	assert.Nil(t, err)
	assert.Equal(t, metrics.EventProcessStatusDiscarded, received)
	eventProcessor.AssertNumberOfCalls(t, "Process", 0)
}

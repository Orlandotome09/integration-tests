package complianceCommandProcessor

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

func TestProcess(t *testing.T) {
	eventProcessor := &mocks.EventProcessor{}

	processor := New(eventProcessor)

	id := uuid.New()
	commmand := &complianceCommand{
		EventType:  "PROFILE_CHANGE",
		EntityID:   id,
		EngineName: "PROFILE",
		Date:       now,
		ParentID:   &id,
		Content:    json.RawMessage(`{"name": "test"}`),
	}
	b, _ := json.Marshal(commmand)

	message := &mocks.Message{}
	message.On("Data").Return(b)
	message.On("PublishTime").Return(now)

	json.Unmarshal(b, &commmand)

	event := &values.Event{
		EngineName:  commmand.EngineName,
		EventType:   commmand.EventType,
		ParentID:    commmand.EntityID,
		EntityID:    commmand.EntityID,
		Date:        commmand.Date,
		Content:     commmand.Content,
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
	msg.On("PublishTime").Return(now)

	received, err := registrationEventProcessor.Process(msg)

	assert.Nil(t, err)
	assert.Equal(t, metrics.EventProcessStatusDiscarded, received)
	eventProcessor.AssertNumberOfCalls(t, "Process", 0)
}

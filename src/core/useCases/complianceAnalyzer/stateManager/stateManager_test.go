package statemanager

import (
	"os"
	"testing"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	stateService *mocks.StateService
	manager      _interfaces.StateManager
)

func TestMain(m *testing.M) {
	stateService = &mocks.StateService{}
	manager = NewStateManager(stateService)
	os.Exit(m.Run())
}

func TestGetOrCreateState_should_get_state_for_event_after_execution_time(t *testing.T) {
	location, _ := time.LoadLocation("UTC")
	eventDate := time.Date(2010, 10, 10, 10, 10, 10, 10, location)
	entityID := uuid.New()
	engineName := values.EngineNameProfile

	state := &entity.State{
		ExecutionTime: time.Date(2000, 10, 10, 10, 10, 10, 10, location),
	}

	stateService.On("Get", entityID).Return(state, true, nil)

	expected := state
	received, shouldIgnore, err := manager.GetOrCreateState(eventDate, entityID, engineName)

	assert.Nil(t, err)
	assert.False(t, shouldIgnore)
	assert.Equal(t, expected, received)
}

func TestGetOrCreateState_eventBeforeState(t *testing.T) {
	location, _ := time.LoadLocation("UTC")
	eventDate := time.Date(2000, 10, 10, 10, 10, 10, 10, location)
	entityID := uuid.New()
	engineName := values.EngineNameProfile

	state := &entity.State{
		RequestDate: time.Date(2010, 10, 10, 10, 10, 10, 10, location),
	}

	stateService.On("Get", entityID).Return(state, true, nil)

	var expected = state
	received, shouldIgnore, err := manager.GetOrCreateState(eventDate, entityID, engineName)

	assert.Nil(t, err)
	assert.True(t, shouldIgnore)
	assert.Equal(t, expected, received)
}

func TestGetOrCreateState_should_create_state_and_send_events(t *testing.T) {
	location, _ := time.LoadLocation("UTC")
	eventDate := time.Date(2000, 10, 10, 10, 10, 10, 10, location)
	entityID := uuid.New()
	engineName := values.EngineNameProfile
	var state *entity.State = nil
	createdState := &entity.State{}

	stateService.On("Get", entityID).Return(state, false, nil)

	stateService.On("Create", entityID, engineName).Return(createdState, nil)

	expected := createdState
	received, shouldIgnore, err := manager.GetOrCreateState(eventDate, entityID, engineName)

	assert.Nil(t, err)
	assert.False(t, shouldIgnore)
	assert.Equal(t, expected, received)
}

func TestGetOrCreateState_should_send_event_with_parentid(t *testing.T) {
	location, _ := time.LoadLocation("UTC")
	eventDate := time.Date(2000, 10, 10, 10, 10, 10, 10, location)
	entityID := uuid.New()
	engineName := values.EngineNameProfile
	var state *entity.State = nil
	createdState := &entity.State{}

	stateService.On("Get", entityID).Return(state, false, nil)

	stateService.On("Create", entityID, engineName).Return(createdState, nil)

	expected := createdState
	received, shouldIgnore, err := manager.GetOrCreateState(eventDate, entityID, engineName)

	assert.Nil(t, err)
	assert.False(t, shouldIgnore)
	assert.Equal(t, expected, received)
}

func TestUpdateState(t *testing.T) {
	state := entity.State{}
	requestDate := time.Now()
	executionTime := time.Now()

	stateService.On("Update", state, requestDate, executionTime).Return(nil)

	err := manager.UpdateState(&state, requestDate, executionTime)

	assert.Nil(t, err)
}

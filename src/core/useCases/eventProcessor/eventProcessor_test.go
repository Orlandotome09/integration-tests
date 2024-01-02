package eventProcessor

import (
	"os"
	"testing"
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	complianceAnalyzer     *mocks.ComplianceAnalyzer
	engineFactory          *mocks.EngineFactory
	eventProcessorInstance interfaces.EventProcessor
)

func TestMain(m *testing.M) {
	complianceAnalyzer = &mocks.ComplianceAnalyzer{}
	engineFactory = &mocks.EngineFactory{}
	eventProcessorInstance = New(complianceAnalyzer, engineFactory)
	os.Exit(m.Run())
}

func TestExecuteForEvent(t *testing.T) {
	entityID := uuid.New()
	location, _ := time.LoadLocation("UTC")
	date := time.Date(2000, 10, 10, 10, 10, 10, 10, location)
	event := &values.Event{
		EngineName: values.EngineNameProfile,
		EventType:  values.EventTypeAccountChanged,
		EntityID:   entityID,
		ParentID:   entityID,
		Date:       date,
	}
	state := entity.State{}
	engine := &mocks.Engine{}

	engineFactory.On("CreateEngine", "PROFILE").Return(engine, nil)
	engine.On("Prepare", entityID).Return(nil)
	complianceAnalyzer.On("RunComplianceAnalysis", engine, entityID, mock.Anything, mock.Anything, false).Return(&state, nil)

	expected := &state
	received, err := eventProcessorInstance.ExecuteForEvent(event)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

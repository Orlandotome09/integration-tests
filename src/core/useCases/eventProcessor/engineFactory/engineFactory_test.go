package enginefactory

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateEngine_profileEngine(t *testing.T) {
	profileEngine := &mocks.Engine{}
	contractEngine := &mocks.Engine{}
	factory := NewEngineFactory(profileEngine, contractEngine)

	engineName := values.EngineNameProfile
	expectedEngine := &mocks.Engine{}

	profileEngine.On("NewInstance").Return(expectedEngine)

	receivedEngine, err := factory.CreateEngine(engineName)

	assert.Nil(t, err)
	assert.Equal(t, expectedEngine, receivedEngine)
}

func TestCreateEngine_contractEngine(t *testing.T) {
	profileEngine := &mocks.Engine{}
	contractEngine := &mocks.Engine{}
	factory := NewEngineFactory(profileEngine, contractEngine)

	engineName := "CONTRACT"
	expectedEngine := &mocks.Engine{}

	contractEngine.On("NewInstance").Return(expectedEngine)

	receivedEngine, err := factory.CreateEngine(engineName)

	assert.Nil(t, err)
	assert.Equal(t, expectedEngine, receivedEngine)
}

func TestCreateEngine_noEngine(t *testing.T) {
	profileEngine := &mocks.Engine{}
	contractEngine := &mocks.Engine{}
	factory := NewEngineFactory(profileEngine, contractEngine)

	engineName := "not engine name"

	receivedEngine, err := factory.CreateEngine(engineName)

	assert.NotNil(t, err)
	assert.Nil(t, receivedEngine)
}

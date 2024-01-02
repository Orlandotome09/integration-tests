package personProcessor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestExecuteForEvent(t *testing.T) {
	personSubEngine := &mocks.SubEngine{}
	newPersonSubEngine := &mocks.SubEngine{}
	complianceAnalyzer := &mocks.ComplianceAnalyzer{}

	profileID := uuid.New()
	person := entity.Person{ProfileID: profileID, EntityID: profileID}
	processor := NewCompliancePersonProcessor(personSubEngine, complianceAnalyzer)

	personSubEngine.On("NewInstance").Return(newPersonSubEngine)
	newPersonSubEngine.On("Prepare", person, "OFFER").Return(nil)
	complianceAnalyzer.On("RunComplianceAnalysis", newPersonSubEngine, profileID, mock.Anything, mock.Anything, false).Return(&entity.State{}, nil)

	state, err := processor.ExecuteForPerson(person, "OFFER")

	assert.Nil(t, err)
	assert.NotNil(t, state)
}

package contractrulesfactory

import (
	"os"
	"testing"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/assert"
)

var (
	documentService     *mocks.DocumentAdapter
	documentFileService *mocks.DocumentFileAdapter
	stateService        *mocks.StateService
	factory             interfaces.ContractRulesFactory
)

func TestMain(m *testing.M) {
	documentService = &mocks.DocumentAdapter{}
	documentFileService = &mocks.DocumentFileAdapter{}
	stateService = &mocks.StateService{}
	factory = New(documentService, documentFileService, stateService)
	os.Exit(m.Run())
}

func TestGetRules(t *testing.T) {
	contract := entity.Contract{}

	expected := []entity.Rule{}

	received := factory.GetRules(contract)

	assert.Equal(t, expected, received)
}

package profileRulesFactory

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	personRules "bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/person"
	"os"
	"testing"

	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"github.com/stretchr/testify/assert"
)

var (
	personComplianceProcessor  *mocks.CompliancePersonProcessor
	personRulesFactoryInstance *mocks.PersonRulesFactory
	factory                    interfaces.ProfileRulesFactory
)

func TestMain(m *testing.M) {
	personComplianceProcessor = &mocks.CompliancePersonProcessor{}
	personRulesFactoryInstance = &mocks.PersonRulesFactory{}
	factory = New(personComplianceProcessor, personRulesFactoryInstance)

	os.Exit(m.Run())
}

func TestGetRules(t *testing.T) {
	ruleSetConfig := &entity.RuleSetConfig{
		ManualBlockParams:        &entity.ManualBlockParams{},
		OwnershipStructureParams: &entity.OwnershipStructureParams{},
		ManualValidationParams:   &entity.ManualValidationParams{},
	}
	profile := &entity.Profile{}
	person := entity.Person{}

	rules := []entity.Rule{
		personRules.NewManualBlockAnalyzer(person),
	}

	personRulesFactoryInstance.On("GetRules", *ruleSetConfig, profile.Person).Return(rules)

	received := factory.GetRules(ruleSetConfig, profile)

	assert.Equal(t, 3, len(received))
}

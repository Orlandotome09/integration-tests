package engines

import (
	"testing"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type contractEngineParams struct {
	ruleValidator            *mocks.RuleValidator
	contractService          *mocks.ContractAdapter
	contractRepository       *mocks.ContractRepository
	contractRulesFactory     *mocks.ContractRulesFactory
	complianceTopicPublisher *mocks.EventPublisher
}

func prepareContractEngine() (contractEngineParams, interfaces.Engine) {
	ruleValidator := &mocks.RuleValidator{}
	contractService := &mocks.ContractAdapter{}
	contractRepository := &mocks.ContractRepository{}
	contractRulesFactory := &mocks.ContractRulesFactory{}
	engine := NewContractEngine(ruleValidator, contractService, contractRepository, contractRulesFactory)
	return contractEngineParams{
		ruleValidator:        ruleValidator,
		contractService:      contractService,
		contractRepository:   contractRepository,
		contractRulesFactory: contractRulesFactory,
	}, engine
}

func TestContractEngine_Prepare(t *testing.T) {
	params, engine := prepareContractEngine()
	profileID := uuid.New()
	entityID := uuid.New()
	contract := &entity.Contract{
		ContractID:           &entityID,
		EstimatedTotalAmount: 0,
		DueTime:              "",
		Installments:         0,
		CorrelationID:        "",
		ProfileID:            &profileID,
		DocumentID:           nil,
	}
	contractRules := []entity.Rule{}
	validationRulesByLimit := []entity.RuleValidatorStep{
		{
			SkipForApproval: false,
			Rules:           contractRules,
		},
	}

	params.contractService.On("Get", &entityID).Return(contract, true, nil)
	params.contractRepository.On("Save", *contract).Return(contract, nil)
	params.contractRepository.On("Get", &entityID).Return(contract, nil)
	params.contractRulesFactory.On("GetRules", *contract).Return(contractRules)
	params.ruleValidator.On("SetRules", validationRulesByLimit).Return()

	err := engine.Prepare(entityID)

	asserted := params.ruleValidator.AssertExpectations(t)

	assert.Nil(t, err)
	assert.True(t, asserted)
}

func TestContractEngineGetName(t *testing.T) {
	_, engine := prepareContractEngine()

	expected := values.EngineNameContract
	received := engine.GetName()

	assert.Equal(t, expected, received)
}

func TestContractEngine_NewInstance(t *testing.T) {
	params, engine := prepareContractEngine()

	newRuleValidator := &mocks.RuleValidator{}

	params.ruleValidator.On("NewInstance").Return(newRuleValidator)

	expected := &contractEngine{
		ruleValidator:        newRuleValidator,
		contractAdapter:      params.contractService,
		contractRepository:   params.contractRepository,
		contractRulesFactory: params.contractRulesFactory,
	}

	received := engine.NewInstance()

	assert.Equal(t, expected, received)
}

package engines

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type contractEngine struct {
	ruleValidator        interfaces.RuleValidator
	contractAdapter      interfaces.ContractAdapter
	contractRepository   interfaces.ContractRepository
	contractRulesFactory interfaces.ContractRulesFactory
}

func NewContractEngine(ruleValidator interfaces.RuleValidator,
	contractAdapter interfaces.ContractAdapter,
	contractRepository interfaces.ContractRepository,
	rulesFactory interfaces.ContractRulesFactory,
) interfaces.Engine {
	return &contractEngine{
		ruleValidator:        ruleValidator,
		contractAdapter:      contractAdapter,
		contractRepository:   contractRepository,
		contractRulesFactory: rulesFactory,
	}
}

func (ref *contractEngine) Prepare(entityID uuid.UUID) error {

	if ref.contractAdapter == nil {
		return errors.New("Contract Service not defined in Contract Engine")
	}

	if ref.contractRepository == nil {
		return errors.New("Contract Repository not defined in Contract Engine")
	}

	if ref.contractRulesFactory == nil {
		return errors.New("Contract Rules Factory not defined in Contract Engine")
	}

	contract, exists, err := ref.contractAdapter.Get(&entityID)
	if err != nil {
		return errors.WithStack(err)
	}

	if !exists {
		return values.NewErrorNotFound("Contract")
	}

	_, err = ref.contractRepository.Save(*contract)
	if err != nil {
		return errors.WithStack(err)
	}

	rules := ref.contractRulesFactory.GetRules(*contract)

	validationRulesByLimit := []entity.RuleValidatorStep{
		{
			StepNumber:      0,
			SkipForApproval: false,
			Rules:           rules,
		},
	}

	ref.ruleValidator.SetRules(validationRulesByLimit)

	return nil
}

func (ref *contractEngine) Validate(state entity.State, override entity.Overrides,
	noCache bool, entityID uuid.UUID, engineName string) (*entity.State, error) {

	return ref.ruleValidator.Validate(state, override, noCache, entityID, engineName)
}

func (ref *contractEngine) PosProcessing(previousState *entity.State, newState *entity.State, entityID uuid.UUID) error {
	return nil
}

func (ref *contractEngine) GetName() string {
	return values.EngineNameContract
}

func (ref *contractEngine) SetCatalog(catalog *entity.CadastralValidationConfig) {}

func (ref *contractEngine) NewInstance() interfaces.Engine {

	return &contractEngine{
		ruleValidator:        ref.ruleValidator.NewInstance(),
		contractAdapter:      ref.contractAdapter,
		contractRepository:   ref.contractRepository,
		contractRulesFactory: ref.contractRulesFactory,
	}
}

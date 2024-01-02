package personRulesConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type personRulesConstructor struct {
	personRulesFactory interfaces.PersonRulesFactory
}

func New(personRulesFactory interfaces.PersonRulesFactory) interfaces.PersonConstructor {
	return &personRulesConstructor{
		personRulesFactory: personRulesFactory}
}

func (ref *personRulesConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.HasCadastralValidationConfig() {
		return nil
	}

	validatorSteps := make([]entity.RuleValidatorStep, 0)

	for _, validationStep := range personWrapper.Person.CadastralValidationConfig.ValidationSteps {
		if validationStep.HasRules() {
			rules := ref.personRulesFactory.GetRules(*validationStep.RulesConfig, personWrapper.Person)

			validatorStep := entity.RuleValidatorStep{
				StepNumber:      validationStep.StepNumber,
				SkipForApproval: validationStep.SkipForApproval,
				Rules:           rules,
			}

			validatorSteps = append(validatorSteps, validatorStep)
		}
	}
	personWrapper.Person.ValidationSteps = validatorSteps

	return nil
}

package rulesConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type profileRulesConstructor struct {
	profileRulesFactory interfaces.ProfileRulesFactory
}

func New(profileRulesFactory interfaces.ProfileRulesFactory) interfaces.ProfileConstructor {
	return &profileRulesConstructor{
		profileRulesFactory: profileRulesFactory}
}

func (ref *profileRulesConstructor) Assemble(profileWrapper *entity.ProfileWrapper) error {

	validationSteps := make([]entity.RuleValidatorStep, 0)

	if profileWrapper.Profile.CadastralValidationConfig == nil {
		return nil
	}

	for _, step := range profileWrapper.Profile.CadastralValidationConfig.ValidationSteps {

		if step.RulesConfig == nil {
			// RuleSet config not configured for step, or it does not exist
			continue
		}
		validationRules := ref.profileRulesFactory.GetRules(step.RulesConfig, &profileWrapper.Profile)

		step := entity.RuleValidatorStep{
			StepNumber:      step.StepNumber,
			SkipForApproval: step.SkipForApproval,
			Rules:           validationRules,
		}

		validationSteps = append(validationSteps, step)

	}

	profileWrapper.Profile.ValidationSteps = validationSteps

	return nil
}

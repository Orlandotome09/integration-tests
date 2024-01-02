package entity

type ValidationStep struct {
	StepNumber      int            `json:"step_number"`
	RulesConfig     *RuleSetConfig `json:"rules_config,omitempty"`
	SkipForApproval bool           `json:"skip_for_approval"`
}

type ValidationSteps []ValidationStep

func (validationStep ValidationStep) HasRules() bool {
	return validationStep.RulesConfig != nil
}

func (validationSteps ValidationSteps) HaveDocumentsValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.IncompleteParams != nil &&
			validationStep.RulesConfig.IncompleteParams.DocumentsRequired != nil {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HaveORCValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() && validationStep.RulesConfig.DOAParams != nil {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HaveAddressValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.IncompleteParams != nil &&
			validationStep.RulesConfig.IncompleteParams.AddressRequired {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HaveWatchlistValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.WatchListParams != nil {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HavePEPValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.PepParams != nil {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HaveBlacklistValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.BlackListParams != nil {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HaveActivityRiskValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.ActivityRiskParams != nil {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HaveBureauValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.BureauParams != nil {
			return true
		}
	}

	return false
}

func (validationSteps ValidationSteps) HaveCAFValidation() bool {
	for _, validationStep := range validationSteps {
		if validationStep.HasRules() &&
			validationStep.RulesConfig.CAFParams != nil {
			return true
		}
	}

	return false
}

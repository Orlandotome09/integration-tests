package entity

type RuleValidatorStep struct {
	StepNumber      int    `json:"step_number"`
	SkipForApproval bool   `json:"skip_for_approval"`
	Rules           []Rule `json:"-"`
}

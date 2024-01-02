package contract

type RuleValidatorStep struct {
	StepNumber      int    `json:"step_number"`
	SkipForApproval bool   `json:"skip_for_approval"`
	Rules           []Rule `json:"rules"`
}

type Rule interface {
	Analyze() ([]RuleResultV2, error)
	Name() string
}

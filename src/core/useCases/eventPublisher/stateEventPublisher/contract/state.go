package contract

import (
	"github.com/google/uuid"
	"time"
)

type State struct {
	EntityID               uuid.UUID              `json:"entity_id"`
	EngineName             string                 `json:"engine_name"`
	Result                 string                 `json:"result"`
	ValidationStepsResults []ValidationStepResult `json:"validation_steps_results"`
	RuleNames              []string               `json:"rule_names"`
	Pending                bool                   `json:"pending"`
	ExecutionTime          time.Time              `json:"execution_time"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type ValidationStepResult struct {
	StepNumber      int            `json:"step_number"`
	Result          string         `json:"result"`
	SkipForApproval bool           `json:"skip_for_approval"`
	RuleResults     []RuleResultV2 `json:"rule_results"`
}

type RuleResultV2 struct {
	RuleSet           string      `json:"rule_set"`
	RuleName          string      `json:"rule_name"`
	Result            string      `json:"result"`
	Pending           bool        `json:"pending"`
	ExpireAt          *time.Time  `json:"expire_at"`
	Metadata          interface{} `json:"metadata"`
	Tags              []string    `json:"tags"`
	Problems          []Problem   `json:"problems"`
	ApprovedDocuments []string    `json:"approved_documents"`
}

type Problem struct {
	Code   string      `json:"code"`
	Detail interface{} `json:"detail"`
}

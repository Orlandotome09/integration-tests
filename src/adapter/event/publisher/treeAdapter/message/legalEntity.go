package message

import "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message/types"

type LegalEntity struct {
	StateNumber      int                    `json:"state_number,omitempty"`
	State            types.StateCode        `json:"state,omitempty"`
	Size             types.CompanySize      `json:"size" binding:"required"`
	BusinessActivity types.BusinessActivity `json:"business_activity" binding:"required"`
	MonthlyIncome    float64                `json:"monthly_income,omitempty"`
	LegalNature      string                 `json:"legal_nature,omitempty"`
}

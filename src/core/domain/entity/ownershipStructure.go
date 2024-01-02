package entity

type OwnershipStructure struct {
	FinalBeneficiariesCount int           `json:"final_beneficiaries_count"`
	ShareholdingSum         float64       `json:"shareholding_sum"`
	Shareholders            []Shareholder `json:"shareholders"`
}

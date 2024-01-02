package contracts

type OwnershipStructureResponse struct {
	LegalEntityID             string        `json:"legal_entity_id"`
	FinalBeneficiariesCounted int           `json:"final_beneficiaries_counted"`
	ShareholdingSum           float64       `json:"shareholding_sum"`
	Shareholders              []Shareholder `json:"shareholders"`
}

type Shareholder struct {
	ParentLegalEntity string       `json:"parent_legal_entity"`
	Shareholding      float64      `json:"shareholding"`
	Role              string       `json:"role"`
	Type              string       `json:"type"`
	Name              string       `json:"name"`
	DocumentNumber    string       `json:"document_number"`
	Nationality       string       `json:"nationality"`
	BirthDate    string       `json:"birth_date"`
	Shareholders Shareholders `json:"shareholders,omitempty"`
}

type Shareholders []Shareholder

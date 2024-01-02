package contract

type EnrichedInformation struct {
	BureauStatus string `json:"status,omitempty"`
	EnrichedIndividual
	EnrichedCompany
	Providers []Provider `json:"providers"`
}

type EnrichedIndividual struct {
	Name      string `json:"name,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
}

type EnrichedCompany struct {
	LegalName          string              `json:"legal_name,omitempty"`
	EconomicActivity   string              `json:"economic_activity,omitempty"`
	OpeningDate        string              `json:"opening_date,omitempty"`
	LegalNature        string              `json:"legal_nature,omitempty"`
	BoardOfDirectors   []Director          `json:"board_of_directors,omitempty"`
	OwnershipStructure *OwnershipStructure `json:"ownership_structure,omitempty"`
}

type Provider struct {
	ProviderName string `json:"provider_name"`
	RequestID    string `json:"request_id,omitempty"`
	Status       string `json:"status,omitempty"`
}

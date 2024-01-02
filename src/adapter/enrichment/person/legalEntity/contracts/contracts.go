package contracts

type LegalEntityResponse struct {
	LegalName        string     `json:"legal_name"`
	OpeningDate      string     `json:"opening_date"`
	Situation        int        `json:"situation"`
	CNAE             string     `json:"cnae"`
	LegalNature      string     `json:"legal_nature"`
	BoardOfDirectors []Director `json:"board_of_directors,omitempty"`
}

type Director struct {
	DocumentNumber string `json:"document_number"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	PersonType     string `json:"person_type"`
}

package enricherContracts

import "github.com/google/uuid"

type EnricherResponse struct {
	Person
	Providers []Provider `json:"providers"`
}

type Person struct {
	EntityID       uuid.UUID            `json:"entity_id"`
	Role           string               `json:"role"`
	Type           string               `json:"type"`
	Name           string               `json:"name"`
	DocumentNumber string               `json:"document_number"`
	Addresses      []Address            `json:"addresses"`
	PhoneNumber    string               `json:"phone_number"`
	Email          string               `json:"email"`
	DocumentFiles  []DocumentFile       `json:"document_files"`
	Individual     *IndividualResponse  `json:"individual,omitempty"`
	Company        *LegalEntityResponse `json:"company,omitempty"`
}

type Provider struct {
	ProviderName string `json:"provider_name"`
	RequestID    string `json:"request_id,omitempty"`
	Status       string `json:"status,omitempty"`
}

type IndividualResponse struct {
	Nationality string `json:"nationality,omitempty"`
	BirthDate   string `json:"birth_date,omitempty"`
	Situation   int    `json:"situation"`
}

type LegalEntityResponse struct {
	BusinessName     string        `json:"business_name"`
	CNAE             string        `json:"cnae"`
	OpeningDate      string        `json:"opening_date"`
	LegalNature      string        `json:"legal_nature"`
	Size             string        `json:"size"`
	Situation        int           `json:"situation"`
	Shareholders     []Shareholder `json:"shareholders,omitempty"`
	BoardOfDirectors []Director    `json:"board_of_directors,omitempty"`
}

type Address struct {
	Street       string `json:"street,omitempty"`
	Number       string `json:"number,omitempty"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Country      string `json:"country,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
}

type Shareholder struct {
	Person
	ParentLegalEntity string        `json:"parent_legal_entity"`
	Shareholding      float64       `json:"shareholding"`
	Nationality       string        `json:"nationality"`
	BirthDate         string        `json:"birth_date"`
	Shareholders      []Shareholder `json:"shareholders,omitempty"`
}

type Director struct {
	Person
}

type DocumentFile struct {
	Type     string `json:"type"`
	FileSide string `json:"file_side"`
	URL      string `json:"url"`
}

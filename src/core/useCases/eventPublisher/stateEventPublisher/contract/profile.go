package contract

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	Person               `json:"person"`
	ProfileID            *uuid.UUID            `json:"profile_id,omitempty"`
	ParentID             *uuid.UUID            `json:"parent_id,omitempty"`
	LegacyID             string                `json:"legacy_id,omitempty"`
	CallbackUrl          string                `json:"callback_url,omitempty"`
	LegalRepresentatives []LegalRepresentative `json:"legal_representatives,omitempty"`
	OwnershipStructure   *OwnershipStructure   `json:"ownership_structure,omitempty"`
	BoardOfDirectors     []Director            `json:"board_of_directors,omitempty"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

type LegalRepresentative struct {
	Person
	LegalRepresentativeID uuid.UUID `json:"legal_representative_id"`
}

type OwnershipStructure struct {
	FinalBeneficiariesCount int           `json:"final_beneficiaries_count"`
	ShareholdingSum         float64       `json:"shareholding_sum"`
	Shareholders            []Shareholder `json:"shareholders"`
}

type Shareholder struct {
	Person           Person     `json:"person"`
	ShareholderID    *uuid.UUID `json:"shareholder_id"`
	ParentID         *uuid.UUID `json:"parent_id"`
	Level            int        `json:"level"`
	OwnershipPercent float64    `json:"ownership_percent"`
}

type Director struct {
	Person     Person    `json:"person"`
	DirectorID uuid.UUID `json:"director_id"`
	Role       string    `json:"role"`
}

func NewProfileFromDomain(profile entity.Profile) Profile {
	return Profile{}
}

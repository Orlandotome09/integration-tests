package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"time"
)

type ProfileResponse struct {
	Person               entity.Person                `json:"person"`
	ProfileID            uuid.UUID                    `json:"profile_id"`
	ParentID             uuid.UUID                    `json:"parent_id,omitempty"`
	LegacyID             string                       `json:"legacy_id,omitempty"`
	CallbackUrl          string                       `json:"callback_url,omitempty"`
	LegalRepresentatives []entity.LegalRepresentative `json:"legal_representatives,omitempty"`
	OwnershipStructure   *entity.OwnershipStructure   `json:"ownership_structure,omitempty"`
	BoardOfDirectors     []entity.Director            `json:"board_of_directors,omitempty"`
	ExpirationDate       *time.Time                   `json:"expiration_date,omitempty"`
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`
}

func (ref ProfileResponse) FromDomain(profile entity.Profile) ProfileResponse {
	ref.Person = profile.Person
	ref.ProfileID = *profile.ProfileID
	if profile.ParentID != nil {
		ref.ParentID = *profile.ParentID
	}
	ref.LegacyID = profile.LegacyID
	ref.CallbackUrl = profile.CallbackUrl
	ref.LegalRepresentatives = profile.LegalRepresentatives
	ref.BoardOfDirectors = profile.BoardOfDirectors
	ref.OwnershipStructure = profile.OwnershipStructure
	ref.ExpirationDate = profile.ExpirationDate
	ref.CreatedAt = profile.CreatedAt
	ref.UpdatedAt = profile.UpdatedAt

	return ref
}

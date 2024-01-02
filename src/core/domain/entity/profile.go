package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Profile struct {
	Person               `json:"person"`
	ProfileID            *uuid.UUID            `json:"profile_id,omitempty"`
	ParentID             *uuid.UUID            `json:"parent_id,omitempty"`
	LegacyID             string                `json:"legacy_id,omitempty"`
	CallbackUrl          string                `json:"callback_url,omitempty"`
	CallbackURLs         []CallbackURL         `json:"callback_urls,omitempty"`
	LegalRepresentatives []LegalRepresentative `json:"legal_representatives,omitempty"`
	OwnershipStructure   *OwnershipStructure   `json:"ownership_structure,omitempty"`
	BoardOfDirectors     []Director            `json:"board_of_directors,omitempty"`
	ExpirationDate       *time.Time            `json:"expiration_date,omitempty"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

type CallbackURL struct {
	URL              string `json:"url"`
	NotificationType string `json:"notification_type"`
}

type ProfileWrapper struct {
	Profile Profile
	Mutex   sync.Mutex
}

func (profile Profile) ShouldGetBoardOfDirectors() bool {

	if profile.Person.PersonType != values.PersonTypeCompany {
		return false
	}

	if profile.Person.CadastralValidationConfig == nil {
		return false
	}

	for _, step := range profile.Person.CadastralValidationConfig.ValidationSteps {
		if step.RulesConfig == nil {
			continue
		}
		if step.RulesConfig.BoardOfDirectorsParams != nil {
			return true
		}
	}

	return false
}

func (profile Profile) ShouldGetLegalRepresentatives() bool {

	if profile.Person.PersonType != values.PersonTypeCompany {
		return false
	}

	if profile.Person.CadastralValidationConfig == nil {
		return false
	}

	for _, step := range profile.Person.CadastralValidationConfig.ValidationSteps {
		if step.RulesConfig == nil {
			continue
		}
		if step.RulesConfig.LegalRepresentativeParams != nil {
			return true
		}
	}

	return false
}

func (profile Profile) ShouldGetOwnershipStructure() bool {

	if profile.Person.PersonType != values.PersonTypeCompany {
		return false
	}

	if profile.Person.CadastralValidationConfig == nil {
		return false
	}

	for _, step := range profile.Person.CadastralValidationConfig.ValidationSteps {
		if step.RulesConfig == nil {
			continue
		}
		if step.RulesConfig.OwnershipStructureParams != nil {
			return true
		}
	}

	return false
}

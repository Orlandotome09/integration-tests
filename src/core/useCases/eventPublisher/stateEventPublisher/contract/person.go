package contract

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type Person struct {
	DocumentNumber         string                  `json:"document_number"`
	Name                   string                  `json:"name"`
	PersonType             string                  `json:"person_type"`
	Email                  string                  `json:"email"`
	PartnerID              string                  `json:"partner_id"`
	OfferType              string                  `json:"offer_type"`
	ProfileID              uuid.UUID               `json:"profile_id"`
	EntityID               uuid.UUID               `json:"entity_id"`
	EntityType             string                  `json:"entity_type"`
	RoleType               string                  `json:"role_type"`
	Individual             *Individual             `json:"individual,omitempty"`
	Company                *Company                `json:"company,omitempty"`
	EnrichedInformation    *EnrichedInformation    `json:"enriched_information,omitempty"`
	BlacklistStatus        *BlacklistStatus        `json:"blacklist_status,omitempty"`
	Watchlist              *Watchlist              `json:"watchlist,omitempty"`
	Addresses              []Address               `json:"addresses,omitempty"`
	Contacts               []Contact               `json:"contacts,omitempty"`
	NotificationRecipients []NotificationRecipient `json:"notification_recipients,omitempty"`
	Documents              []Document              `json:"documents,omitempty"`
	DocumentFiles          []DocumentFile          `json:"document_files,omitempty"`
	Overrides              []Override              `json:"overrides,omitempty"`
	ValidationSteps        []RuleValidatorStep     `json:"validation_steps,omitempty"`
	Catalog                *Catalog                `json:"catalog,omitempty"`
}

func NewPersonFromDomain(person entity.Person) Person {
	return Person{}
}

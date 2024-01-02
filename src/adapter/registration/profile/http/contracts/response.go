package contracts

import (
	entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"time"

	"github.com/google/uuid"
)

type ProfileResponse struct {
	ProfileID      uuid.UUID     `json:"profile_id"`
	PartnerID      string        `json:"partner_id"`
	DocumentNumber string        `json:"document_number"`
	ParentID       *uuid.UUID    `json:"parent_id,omitempty"`
	Name           string        `json:"name,omitempty"`
	LegacyID       string        `json:"legacy_id,omitempty"`
	OfferType      string        `json:"offer_type,omitempty"`
	RoleType       string        `json:"role_type,omitempty"`
	ProfileType    string        `json:"profile_type,omitempty"`
	CallbackUrl    string        `json:"callback_url,omitempty"`
	CallbackURLs   []CallbackURL `json:"callback_urls,omitempty"`
	Individual     *Individual   `json:"individual,omitempty"`
	Company        *Company      `json:"company,omitempty"`
	Email          string        `json:"email,omitempty"`
	ExpirationDate *time.Time    `json:"expiration_date,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type CallbackURL struct {
	URL              string `json:"url"`
	NotificationType string `json:"notification_type"`
}

func (ref *ProfileResponse) ToDomain() *entity.Profile {
	if ref == nil {
		return nil
	}

	var callbackURLs []entity.CallbackURL

	for _, callbackURL := range ref.CallbackURLs {
		callbackURLs = append(callbackURLs, entity.CallbackURL{
			URL:              callbackURL.URL,
			NotificationType: callbackURL.NotificationType,
		})
	}

	return &entity.Profile{
		ProfileID: &ref.ProfileID,
		Person: entity.Person{
			DocumentNumber: ref.DocumentNumber,
			Name:           ref.Name,
			PersonType:     ref.ProfileType,
			Email:          ref.Email,
			PartnerID:      ref.PartnerID,
			ProfileID:      ref.ProfileID,
			EntityID:       ref.ProfileID,
			EntityType:     values.EntityTypeProfile,
			OfferType:      ref.OfferType,
			Individual:     ref.Individual.ToDomain(),
			Company:        ref.Company.ToDomain(),
			RoleType:       ref.RoleType,
		},
		ParentID:       ref.ParentID,
		LegacyID:       ref.LegacyID,
		CallbackUrl:    ref.CallbackUrl,
		CallbackURLs:   callbackURLs,
		ExpirationDate: ref.ExpirationDate,
		CreatedAt:      ref.CreatedAt,
		UpdatedAt:      ref.UpdatedAt,
	}
}

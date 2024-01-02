package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"time"
)

type ProfileRequest struct {
	PartnerID      string             `json:"partner_id"`
	DocumentNumber string             `json:"document_number"`
	ParentID       *uuid.UUID         `json:"parent_id,omitempty"`
	Name           string             `json:"name"`
	LegacyID       string             `json:"legacy_id,omitempty"`
	OfferType      string             `json:"offer_type"`
	RoleType       values2.RoleType   `json:"role_type"`
	ProfileType    values2.PersonType `json:"profile_type"`
	CallbackUrl    string             `json:"callback_url,omitempty"`
	CallbackURLs   []CallbackURL      `json:"callback_urls,omitempty"`
	Email          string             `json:"email,omitempty"`
	ExpirationDate *time.Time         `json:"expiration_date,omitempty"`
	Individual     *Individual        `json:"individual,omitempty"`
	Company        *Company           `json:"company,omitempty"`
}

func NewProfileRequestFromDomain(profile *entity.Profile) *ProfileRequest {
	if profile == nil {
		return nil
	}

	var callbackURLs []CallbackURL

	for _, callbackURL := range profile.CallbackURLs {
		callbackURLs = append(callbackURLs, CallbackURL{
			URL:              callbackURL.URL,
			NotificationType: callbackURL.NotificationType,
		})
	}
	return &ProfileRequest{
		PartnerID:      profile.PartnerID,
		DocumentNumber: profile.DocumentNumber,
		ParentID:       profile.ParentID,
		Name:           profile.Name,
		LegacyID:       profile.LegacyID,
		OfferType:      profile.OfferType,
		RoleType:       profile.RoleType,
		ProfileType:    profile.Person.PersonType,
		CallbackUrl:    profile.CallbackUrl,
		CallbackURLs:   callbackURLs,
		Email:          profile.Email,
		ExpirationDate: profile.ExpirationDate,
		Individual:     NewIndividualFromDomain(profile.Individual),
		Company:        NewCompanyFromDomain(profile.Company),
	}
}

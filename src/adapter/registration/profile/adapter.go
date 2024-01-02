package profile

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	profileHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile/http/contracts"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type profileAdapter struct {
	interfaces.ProfileAdapter
	profileClient profileHttpClient.ProfileHttpClient
}

func NewProfileService(profileClient profileHttpClient.ProfileHttpClient) interfaces.ProfileAdapter {
	return &profileAdapter{
		profileClient: profileClient,
	}
}

func (ref *profileAdapter) Get(profileID uuid.UUID) (*entity.Profile, error) {
	resp, err := ref.profileClient.Get(profileID)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.ToDomain(), nil
}

func (ref *profileAdapter) Create(profile *entity.Profile) (*entity.Profile, error) {
	profileRequest := contracts.NewProfileRequestFromDomain(profile)

	resp, err := ref.profileClient.Create(*profileRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.ToDomain(), nil
}

func (ref *profileAdapter) FindByDocumentNumber(roleType values.RoleType, documentNumber string, partnerID string, parentID *uuid.UUID) (*entity.Profile, error) {
	resp, err := ref.profileClient.FindByDocumentNumber(roleType, documentNumber, partnerID, parentID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.ToDomain(), nil
}

package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type ProfileAdapter interface {
	Get(profileID uuid.UUID) (currentProfile *entity.Profile, err error)
	Create(profile *entity.Profile) (savedProfile *entity.Profile, err error)
	FindByDocumentNumber(roleType values.RoleType, documentNumber string, partnerID string, parentID *uuid.UUID) (*entity.Profile, error)
}

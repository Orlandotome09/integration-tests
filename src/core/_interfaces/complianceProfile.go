package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type ComplianceProfileService interface {
	Get(profileID uuid.UUID) (*entity.Profile, error)
}

type ComplianceProfileRepository interface {
	Save(profile entity.Profile) (*entity.Profile, error)
	Get(profileID uuid.UUID) (*entity.Profile, error)
}

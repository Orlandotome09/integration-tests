package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type ProfileFactory interface {
	Build(profileID uuid.UUID) (*entity.Profile, error)
}

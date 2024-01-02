package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type LegalRepresentativeAdapter interface {
	Get(id uuid.UUID) (*entity.LegalRepresentative, error)
	Search(profileID uuid.UUID) ([]entity.LegalRepresentative, error)
}

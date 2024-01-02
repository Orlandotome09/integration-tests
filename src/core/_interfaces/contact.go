package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type ContactAdapter interface {
	Get(id uuid.UUID) (*entity.Contact, error)
	Search(profileID string) ([]entity.Contact, error)
}

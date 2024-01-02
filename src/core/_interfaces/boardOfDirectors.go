package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type BoardOfDirectorsAdapter interface {
	Search(profileID uuid.UUID) ([]entity.Director, error)
}

package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type PersonRepository interface {
	Save(person entity.Person) (*entity.Person, error)
	Get(personID uuid.UUID) (*entity.Person, error)
}

package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type PersonFactory interface {
	Build(person entity.Person) (*entity.Person, error)
}

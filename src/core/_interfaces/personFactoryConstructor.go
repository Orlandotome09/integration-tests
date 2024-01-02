package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type PersonConstructor interface {
	Assemble(person *entity.PersonWrapper) error
}

package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ProfileConstructor interface {
	Assemble(profile *entity.ProfileWrapper) error
}

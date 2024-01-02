package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type AddressAdapter interface {
	Get(id string) (*entity.Address, error)
	Search(profileID string) ([]entity.Address, error)
}

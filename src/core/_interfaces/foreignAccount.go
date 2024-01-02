package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ForeignAccountAdapter interface {
	Get(foreignAccountID string) (*entity.ForeignAccount, error)
}

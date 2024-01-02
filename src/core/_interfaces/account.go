package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type AccountAdapter interface {
	GetByID(accountID string) (*entity.Account, error)
	FindByProfileID(profileID string) ([]entity.Account, error)
	CreateInternal(entityId string, bankCode string) (*entity.Account, error)
}

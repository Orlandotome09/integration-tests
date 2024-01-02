package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type BlacklistAdapter interface {
	Search(documentNumber, partnerId string) (*entity.BlacklistStatus, bool, error)
}

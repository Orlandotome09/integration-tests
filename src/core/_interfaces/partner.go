package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type PartnerAdapter interface {
	GetActive(partnerID string) (*entity.Partner, error)
}

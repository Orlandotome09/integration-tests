package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type EconomicActivityRepository interface {
	Get(code string) (record *entity.EconomicActivity, exists bool, err error)
}

type EconomicActivityService interface {
	Get(code string) (record *entity.EconomicActivity, exists bool, err error)
}

package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type BureauService interface {
	GetBureauStatus(person entity.Person) (*entity.EnrichedInformation, error)
}

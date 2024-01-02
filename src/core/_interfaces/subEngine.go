package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type SubEngine interface {
	Prepare(person entity.Person, offerType string) error
	NewInstance() SubEngine
	ComplianceValidator
}

type SubEngineFactory interface {
	CreateSubEngine(subEngineName string) (SubEngine, error)
}

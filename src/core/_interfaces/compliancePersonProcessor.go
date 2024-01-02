package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type CompliancePersonProcessor interface {
	ExecuteForPerson(person entity2.Person, offer string) (*entity2.State, error)
}

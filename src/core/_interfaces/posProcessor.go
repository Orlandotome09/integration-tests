package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type PosProcessor interface {
	SendToTreeAdapter(profile entity2.Profile, catalogConfig entity2.ProductConfig) error
	SendToLimit(profile entity2.Profile, state entity2.State, catalogConfig entity2.ProductConfig) error
	CreateInternalAccount(entityID uuid.UUID) error
}

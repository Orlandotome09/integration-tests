package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"context"
)

type TemisConfigAdapter interface {
	GetCadastralValidationConfig(ctx context.Context, personType, roleType, offerType, partnerID string) (*entity.CadastralValidationConfig, error)
}

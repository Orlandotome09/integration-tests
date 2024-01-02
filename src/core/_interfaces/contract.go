package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type ContractAdapter interface {
	Get(contractId *uuid.UUID) (*entity.Contract, bool, error)
}

type ContractRepository interface {
	Get(id uuid.UUID) (*entity.Contract, error)
	Save(contract entity.Contract) (*entity.Contract, error)
}

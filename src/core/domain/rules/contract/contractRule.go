package contractRule

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ContractRule struct {
	entity.Rule
	Contract entity.Contract
}

package contract

import (
	contractClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type contractAdapter struct {
	interfaces.ContractAdapter
	contractClient contractClient.ContractClient
}

func NewContractAdapter(contractClient contractClient.ContractClient) interfaces.ContractAdapter {
	return &contractAdapter{contractClient: contractClient}
}

func (ref *contractAdapter) Get(contractId *uuid.UUID) (*entity.Contract, bool, error) {

	contractResponse, exists, err := ref.contractClient.Get(contractId)
	if err != nil {
		return nil, false, err
	}

	if exists == false {
		return nil, false, nil
	}

	contract, err := contractResponse.ToDomain()
	if err != nil {
		return nil, false, err
	}

	return contract, true, nil
}

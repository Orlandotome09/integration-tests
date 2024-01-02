package contractrulesfactory

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type contractRulesFactory struct {
	documentService     interfaces.DocumentAdapter
	documentFileService interfaces.DocumentFileAdapter
	stateService        interfaces.StateService
}

func New(documentService interfaces.DocumentAdapter,
	documentFileService interfaces.DocumentFileAdapter,
	stateService interfaces.StateService) interfaces.ContractRulesFactory {
	return &contractRulesFactory{
		documentService:     documentService,
		documentFileService: documentFileService,
		stateService:        stateService,
	}
}

func (ref *contractRulesFactory) GetRules(contract entity.Contract) []entity.Rule {
	return []entity.Rule{}
}

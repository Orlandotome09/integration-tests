package ownershipStructure

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/idgenerator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type ownershipStructureService struct {
	httpClient  ownershipStructureClient.OwnershipStructureClient
	translator  translator.OwnershipStructureTranslator
	idgenerator idgenerator.IdGenerator
}

func New(httpClient ownershipStructureClient.OwnershipStructureClient,
	translator translator.OwnershipStructureTranslator,
	idgenerator idgenerator.IdGenerator) interfaces.OwnershipStructureAdapter {
	return &ownershipStructureService{
		httpClient:  httpClient,
		translator:  translator,
		idgenerator: idgenerator,
	}
}

func (ref *ownershipStructureService) Get(legalEntityID, offerType, partnerID string) (*entity.OwnershipStructure, error) {
	response, err := ref.httpClient.GetOwnershipStructure(legalEntityID, offerType, partnerID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if response == nil {
		return nil, nil
	}

	ownershipStructure, err := ref.translator.Translate(*response)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i, shareholder := range ownershipStructure.Shareholders {
		id := ref.idgenerator.Generate(response.LegalEntityID, shareholder.Person.DocumentNumber)
		ownershipStructure.Shareholders[i].ShareholderID = &id
		ownershipStructure.Shareholders[i].Person.EntityID = id
	}

	return ownershipStructure, nil
}

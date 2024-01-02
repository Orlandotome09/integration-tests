package ownershipStructure

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ownershipStructureAdapter struct {
	httpClient ownershipStructureHttpClient.OwnershipStructureClient
	translator translator.OwnershipStructureTranslator
}

func New(httpClient ownershipStructureHttpClient.OwnershipStructureClient,
	translator translator.OwnershipStructureTranslator) interfaces.OwnershipStructureAdapter {
	return &ownershipStructureAdapter{
		httpClient: httpClient,
		translator: translator,
	}
}

func (ref *ownershipStructureAdapter) Get(profileID, offerType, partnerID string) (*entity.OwnershipStructure, error) {
	response, err := ref.httpClient.GetOwnershipStructureForProfile(profileID)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}

	return ref.translator.Translate(*response)
}

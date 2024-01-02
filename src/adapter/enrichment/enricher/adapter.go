package enricherAdapter

import (
	enricherClient "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http/contracts"
	enricherTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type enricherAdapter struct {
	httpClient enricherClient.EnricherClient
	translator enricherTranslator.EnricherTranslator
}

func New(httpClient enricherClient.EnricherClient,
	translator enricherTranslator.EnricherTranslator) interfaces.EnricherAdapter {
	return &enricherAdapter{
		httpClient: httpClient,
		translator: translator,
	}
}

func (ref *enricherAdapter) GetEnrichedPerson(documentNumber, profileID, personType, offerType, partnerID, roleTye string) (*entity.EnrichedInformation, error) {
	request := enricherContracts.EnricherRequest{
		ProfileID:  profileID,
		PersonType: personType,
		OfferType:  offerType,
		PartnerID:  partnerID,
		RoleType:   roleTye,
	}
	response, err := ref.httpClient.GetEnrichedPerson(request, documentNumber)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if response == nil {
		return nil, nil
	}

	enrichedInformation := ref.translator.Translate(*response)

	return enrichedInformation, nil
}

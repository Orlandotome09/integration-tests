package personAdapter

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/_interfacesEnrichment"
	enrichedIndividualTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/individual/translator"
	enrichedLegalEntityTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/legalEntity/translator"
	"bitbucket.org/bexstech/temis-compliance/src/core"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

const (
	GetIndividualPath  = "/individual/"
	GetLegalEntityPath = "/legal-entity/"
)

type personAdapter struct {
	enrichedHttpClient            adapter.HttpClient
	enrichedIndividualTranslator  _interfacesEnrichment.EnrichedInformationTranslator
	enrichedLegalEntityTranslator _interfacesEnrichment.EnrichedInformationTranslator
}

func NewPersonAdapter(enrichedHttpClient adapter.HttpClient) interfaces.BureauService {
	return &personAdapter{
		enrichedHttpClient:            enrichedHttpClient,
		enrichedIndividualTranslator:  enrichedIndividualTranslator.New(),
		enrichedLegalEntityTranslator: enrichedLegalEntityTranslator.New(),
	}

}

func (ref *personAdapter) GetBureauStatus(person entity.Person) (*entity.EnrichedInformation, error) {
	if person.IsCompany() {
		return ref.getBureauStatus(person, GetLegalEntityPath, ref.enrichedLegalEntityTranslator)
	} else {
		return ref.getBureauStatus(person, GetIndividualPath, ref.enrichedIndividualTranslator)
	}
}

func (ref *personAdapter) getBureauStatus(person entity.Person, path string, translator _interfacesEnrichment.EnrichedInformationTranslator) (*entity.EnrichedInformation, error) {
	documentNumber := core.NormalizeDocument(person.DocumentNumber)
	headers := map[string]string{"Offer-Type": person.OfferType, "Partner-Id": person.PartnerID}
	fullPath := path + documentNumber
	response, err := ref.enrichedHttpClient.Get(fullPath, "", headers)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if response == nil {
		return nil, nil
	}
	return translator.Translate(documentNumber, person.ProfileID, response)
}

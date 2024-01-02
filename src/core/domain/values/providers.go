package values

type ProviderName string

const (
	IndividualBureauEnricher  ProviderName = "INDIVIDUAL_BUREAU_ENRICHER"
	LegalEntityBureauEnricher ProviderName = "LEGAL_ENTITY_BUREAU_ENRICHER"
	CAFEnricher               ProviderName = "CAF_ENRICHER"
)

func (provider ProviderName) String() string {
	return string(provider)
}

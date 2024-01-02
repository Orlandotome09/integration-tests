package _init

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/auth"
	enricherAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher"
	enricherClient "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http"
	enricherTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/translator"
	personAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address"
	addressClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http"
	registrationAddressesTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact"
	contactClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http"
	contactTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document"
	documentClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http"
	documentTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile"
	documentFileClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http"
	documentFileTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient"
	notificationRecipientClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http"
	notificationRecipientTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/translator"
	restrictivelists "bitbucket.org/bexstech/temis-compliance/src/adapter/restrictiveLists"
	restrictiveListsHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/restrictiveLists/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/temisConfig"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/temisConfig/http"
	watchlistService "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist"
	watchlistClient "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http"
	watchlistTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/addressConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/blacklistConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/bureauConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/cadastralValidationConfigConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/contactConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/documentsConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/enrichmentConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/notificationRecipientConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/pepInformationConstructor"
	personRulesConstructor "bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/rulesConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/personFactory/constructors/watchlistConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/override"
	"bitbucket.org/bexstech/temis-compliance/src/infra"
)

func buildPersonFactory() interfaces.PersonFactory {
	cacheInstance := infra.NewCache()

	bureauClient := adapter.NewHttpClient(webClient, temisEnrichmentHost)
	enrichmentClient := enricherClient.New(webClient, temisEnrichmentHost)
	watchlistClientInstance := watchlistClient.New(webClient, complianceExtHost, nil)
	restrictiveListsClient := restrictiveListsHttpClient.NewRestrictiveListHttpClient(webClient, temisRestrictiveListsHost)
	temisConfigClient := http.NewTemisConfigHttpClient(webClient, temisConfigHost, cacheInstance)

	if environment != "local" {
		complianceAuthRepository := auth.NewAuthRepository("https://bexs.auth0.com", auth.OAuth{
			ClientId:     temisClientId,
			ClientSecret: temisClientSecret,
			Audience:     complianceTokenAudience,
			GrantType:    grantType,
		})
		watchlistClientInstance = watchlistClient.New(webClient, complianceExtHost, complianceAuthRepository)
	}

	addressClientInstance := addressClient.New(webClient, temisRegistrationHost)
	contactClientInstance := contactClient.New(webClient, temisRegistrationHost)
	notificationRecipientInstance := notificationRecipientClient.New(webClient, temisRegistrationHost)
	documentClientInstance := documentClient.New(webClient, temisRegistrationHost)
	documentFileClientInstance := documentFileClient.New(webClient, temisRegistrationHost)

	addressService := address.NewAddressAdapter(addressClientInstance, registrationAddressesTranslator.New())
	bureauServiceInstance := personAdapter.NewPersonAdapter(bureauClient)
	enrichmentAdapterInstance := enricherAdapter.New(enrichmentClient, enricherTranslator.New())
	restrictiveListsAdapter := restrictivelists.NewRestrictiveListsAdapter(restrictiveListsClient)
	watchlistServiceInstance := watchlistService.New(watchlistClientInstance, watchlistTranslator.New())
	overrideService := override.NewOverrideService(overrideRepository)
	documentService := document.NewDocumentAdapter(documentClientInstance, documentTranslator.New())
	documentFileService := documentFile.NewDocumentFileAdapter(documentFileClientInstance, documentFileTranslator.New())
	contactAdapter := contact.NewContactAdapter(contactClientInstance, contactTranslator.New())
	notificationRecipientAdapter := notificationRecipient.NewNotificationRecipientAdapter(notificationRecipientInstance, notificationRecipientTranslator.New())
	temisConfigAdapter := temisConfig.NewTemisConfigAdapter(temisConfigClient)

	return personFactory.New(overrideService,
		cadastralValidationConfigConstructor.New(temisConfigAdapter),
		personRulesConstructor.New(buildPersonRulesFactory()),
		bureauConstructor.New(bureauServiceInstance),
		[]interfaces.PersonConstructor{
			watchlistConstructor.New(watchlistServiceInstance),
			blacklistConstructor.New(restrictiveListsAdapter),
			addressConstructor.New(addressService),
			documentsConstructor.New(documentService, documentFileService),
			enrichmentConstructor.New(enrichmentAdapterInstance),
			contactConstructor.New(contactAdapter),
			notificationRecipientConstructor.New(notificationRecipientAdapter),
			pepInformationConstructor.New(restrictiveListsAdapter),
		},
	)

}

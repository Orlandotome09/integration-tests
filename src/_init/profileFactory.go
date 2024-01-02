package _init

import (
	ownershipStructureAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure"
	ownershipStructureClient "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/idgenerator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors"
	boardOfDirectorsClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http"
	boardOfDirectorsTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative"
	legalRepresentativeClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http"
	legalRepresentativeTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/translator"
	registrationOwnershiptStructure "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure"
	registrationOwnershipStructureClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure/http"
	registrationTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile"
	profileHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/profileFactory"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/profileFactory/constructors/boardOfDirectorsConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/profileFactory/constructors/legalRepresentativesConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/profileFactory/constructors/ownershipStructureConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/service/profileFactory/constructors/rulesConstructor"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/ownershipStructure"
	"net/http"
	"time"
)

func buildProfileFactory() interfaces.ProfileFactory {
	legalRepresentativeClientInstance := legalRepresentativeClient.New(webClient, temisRegistrationHost)
	enrichedClient := ownershipStructureClient.New(&http.Client{Timeout: 1 * time.Minute}, temisEnrichmentHost)
	manuallyFilledClient := registrationOwnershipStructureClient.New(webClient, temisRegistrationHost)
	boardOfDirectorsClientInstance := boardOfDirectorsClient.New(webClient, temisRegistrationHost)
	profileClient := profileHttpClient.NewProfileHttpClient(webClient, temisRegistrationHost)

	enrichmentOwnershipAdapter := ownershipStructureAdapter.New(enrichedClient, translator.New(), idgenerator.NewIdGenerator())
	registrationOwnershipAdapter := registrationOwnershiptStructure.New(manuallyFilledClient, registrationTranslator.New())

	ownershipStructureService := ownershipStructure.New(enrichmentOwnershipAdapter, registrationOwnershipAdapter)
	boardOfDirectorsService := boardOfDirectors.New(boardOfDirectorsClientInstance, boardOfDirectorsTranslator.New())
	legalRepresentativeService := legalRepresentative.New(legalRepresentativeClientInstance, legalRepresentativeTranslator.New())

	profileAdapter := profile.NewProfileService(profileClient)

	return profileFactory.New(profileAdapter, buildPersonFactory(), rulesConstructor.New(buildProfileRulesFactory()),
		[]interfaces.ProfileConstructor{
			ownershipStructureConstructor.New(ownershipStructureService),
			legalRepresentativesConstructor.New(legalRepresentativeService),
			boardOfDirectorsConstructor.New(boardOfDirectorsService),
		})
}

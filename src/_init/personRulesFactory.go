package _init

import (
	"context"

	"bitbucket.org/bexstech/temis-compliance/src/adapter"
	doaAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/doa"
	doaTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/doa/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/file"
	fileHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/file/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	personRulesFactory "bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/person/factory"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/economicActivity"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventPublisher/complianceCommandPublisher"
	"bitbucket.org/bexstech/temis-compliance/src/infra/queues/pubsubPublisher"
	"github.com/go-playground/validator/v10"
)

func buildPersonRulesFactory() interfaces.PersonRulesFactory {
	validate := validator.New()
	ctx := context.Background()

	// complianceCommandPublisher
	complianceCommandsTopicPublisher := pubsubPublisher.New(ctx, bexsPubsubClient, complianceCommandTopic)
	complianceCommandsPublisher := complianceCommandPublisher.New(complianceCommandsTopicPublisher)

	fileClient := fileHttpClient.New(webClient, temisRegistrationHost)
	fileService := file.NewFileAdapter(fileClient)

	doaHttpClient := adapter.NewHttpClient(webClient, doaHost)
	doaAdapterInstance := doaAdapter.NewDOAAdapter(doaHttpClient, doaTranslator.New(temisComplianceHost))

	doaService := doa.NewDOAService(validate, doaRepository, complianceCommandsPublisher, doaAdapterInstance)

	economicalActivityService := economicActivity.NewEconomicActivityService(validate, economicActivityRepository)

	return personRulesFactory.New(fileService, doaService, economicalActivityService)
}

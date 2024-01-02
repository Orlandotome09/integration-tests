package _init

import (
	"context"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"bitbucket.org/bexstech/temis-compliance/src/adapter"
	doaAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/doa"
	doaTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/doa/translator"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/complianceProfile"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventProcessor"
	enginefactory "bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventProcessor/engineFactory"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventPublisher/complianceCommandPublisher"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/offer"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/override"
	"bitbucket.org/bexstech/temis-compliance/src/infra/metrics"
	"bitbucket.org/bexstech/temis-compliance/src/infra/queues/pubsubPublisher"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/authentication"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/doaApi"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/loCallbackApi"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/logger"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/offerApi"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/overrideApi"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/profileApi"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/stateApi"
)

func StartComplianceApis(ctx context.Context, wg *sync.WaitGroup, external bool, db *gorm.DB) {
	defer wg.Done()

	validate := validator.New()

	// PubSub
	complianceCommandTopicPublisher := pubsubPublisher.New(ctx, bexsPubsubClient, complianceCommandTopic)
	complianceCommandsPublisher := complianceCommandPublisher.New(complianceCommandTopicPublisher)

	// httpclient
	doaHttpClient := adapter.NewHttpClient(webClient, doaHost)

	// Services
	offerService := offer.NewOfferService(validate, offerRepository)
	doaAdapterInstance := doaAdapter.NewDOAAdapter(doaHttpClient, doaTranslator.New(temisComplianceHost))
	doaService := doa.NewDOAService(validate, doaRepository, complianceCommandsPublisher, doaAdapterInstance)
	overrideService := override.NewOverrideService(overrideRepository)
	complianceProfileService := complianceProfile.NewComplianceProfileService(profileRepository, personRepository)

	// CnC--------------------------
	engineFactory := enginefactory.NewEngineFactory(buildProfileEngine(), buildContractEngine())

	cncInstance := eventProcessor.New(buildComplianceAnalyzer(), engineFactory)

	// Starting API
	ginServer := presentation.NewServerHttpGin(PrettyLog())

	ginRouterGroupAlwaysOpened := ginServer.GetGinRouterGroup(getRelativePath(external))

	loCallbackApi.RegisterLoCallbackApi(ginRouterGroupAlwaysOpened)
	presentation.RegisterInfraApi(ginRouterGroupAlwaysOpened, diagnosticMode, db)

	if environment != "prd" {
		ginRouterGroupAlwaysOpened.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	ginRouterGroup := ginServer.GetGinRouterGroup(getRelativePath(external))
	ginRouterGroup.Use(logger.Middleware)
	ginRouterGroup.Use(metrics.GetMiddleware())
	ginRouterGroup.Use(gin.Recovery())

	if external {
		ginRouterGroup.Use(authentication.Middleware(audience))
	}

	doaApi.RegisterDOACallback(ginRouterGroup, doaService)
	offerApi.RegisterOfferApi(ginRouterGroup, offerService)
	stateApi.RegisterStateApi(ginRouterGroup, buildStateService(), cncInstance)
	overrideApi.RegisterOverrideApi(ginRouterGroup, overrideService, cncInstance)
	profileApi.RegisterProfileApi(ginRouterGroup, complianceProfileService, cncInstance)

	err := ginServer.StartServer(internalPortNumber)
	if err != nil {
		logrus.Fatal(err)
	}
}

func getRelativePath(external bool) string {
	if external {
		return "/compliance"
	}
	return "/compliance-int"
}

func PrettyLog() bool {
	return strings.ToUpper(pretty) == "TRUE"
}

package _init

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/infra"
	"bitbucket.org/bexstech/temis-compliance/src/repository/contractRepo"
	"bitbucket.org/bexstech/temis-compliance/src/repository/doaResultRepo"
	economicalActivityRepo "bitbucket.org/bexstech/temis-compliance/src/repository/economicActivityRepo"
	"bitbucket.org/bexstech/temis-compliance/src/repository/errorHandler"
	"bitbucket.org/bexstech/temis-compliance/src/repository/offerRepo"
	"bitbucket.org/bexstech/temis-compliance/src/repository/overrideRepo"
	"bitbucket.org/bexstech/temis-compliance/src/repository/personRepo"
	"bitbucket.org/bexstech/temis-compliance/src/repository/profileRepo"
	"bitbucket.org/bexstech/temis-compliance/src/repository/stateRepo"
)

var (
	// Hosts
	temisRegistrationHost     = os.Getenv("TEMIS_HOST")
	temisComplianceHost       = os.Getenv("TEMIS_COMPLIANCE_HOST")
	temisEnrichmentHost       = os.Getenv("TEMIS_ENRICHMENT_HOST")
	temisRestrictiveListsHost = os.Getenv("TEMIS_RESTRICTIVE_LISTS_HOST")
	doaHost                   = os.Getenv("DOA_HOST")
	temisConfigHost           = os.Getenv("TEMIS_CONFIG_HOST")

	// Pubsub
	pubsubProjectID                         = os.Getenv("PUBSUB_PROJECT_ID")
	bexDigitalProjectID                     = os.Getenv("BEXS_DIGITAL_PROJECT_ID")
	complianceCommandTopic                  = os.Getenv("COMPLIANCE_COMMAND_TOPIC_ID")
	complianceCommanndSubscription          = os.Getenv("COMPLIANCE_COMMAND_SUBSCRIPTION_ID")
	registrationEventsSubscription          = os.Getenv("REGISTRATION_EVENTS_SUBSCRIPTION_ID")
	complianceSubscriptionToEnrichmentTopic = os.Getenv("COMPLIANCE_SUBSCRIPTION_TO_ENRICHMENT_TOPIC")
	treeAdapterTopic                        = os.Getenv("TREE_ADAPTER_TOPIC_ID")
	limitTopic                              = os.Getenv("LIMIT_TOPIC_ID")
	stateEventsTopic                        = os.Getenv("STATE_EVENTS_TOPIC_ID")

	// External Compliance
	temisClientId           = os.Getenv("TEMIS_CLIENT_ID")
	temisClientSecret       = os.Getenv("TEMIS_CLIENT_SECRET")
	complianceTokenAudience = os.Getenv("COMPLIANCE_TOKEN_AUDIENCE")
	complianceExtHost       = os.Getenv("COMPLIANCE_EXT_HOST")

	diagnosticMode = strings.ToUpper(os.Getenv("DIAGNOSTIC_MODE")) == "TRUE"
	environment    = os.Getenv("ENVIRONMENT")
	audience       = os.Getenv("TOKEN_AUDIENCE")

	pretty = os.Getenv("PRETTY")
)

var (
	timeNowGenerator  = func() time.Time { return time.Now() }
	randomIDGenerator = func() uuid.UUID { return uuid.New() }

	// Pubsub clients
	bexsPubsubClient        *pubsub.Client
	bexsDigitalPubsubClient *pubsub.Client

	// Http
	webClient *http.Client
)

// Repositories
var (
	doaRepository              interfaces.DOAResultRepository
	economicActivityRepository interfaces.EconomicActivityRepository
	offerRepository            interfaces.OfferRepository
	overrideRepository         interfaces.OverrideRepository
	personRepository           interfaces.PersonRepository
	profileRepository          interfaces.ComplianceProfileRepository
	stateRepository            interfaces.StateRepository
	contractRepository         interfaces.ContractRepository
)

func InitPubsubClients(ctx context.Context) {
	bexsPubsubClient = infra.OpenPubSubConnection(ctx, pubsubProjectID)
	if bexsPubsubClient == nil {
		logrus.Panicf("error connecting to pubsub project: %v", bexsPubsubClient)
	}
	bexsDigitalPubsubClient = infra.OpenPubSubConnection(ctx, bexDigitalProjectID)
	if bexsDigitalPubsubClient == nil {
		logrus.Panicf("error connecting to pubsub project: %v", bexsDigitalPubsubClient)
	}
}

func InitWebClients() {
	webClient = &http.Client{Timeout: 20 * time.Second}
}

func InitRepositories(db *gorm.DB) {
	dbErrorHandler := errorHandler.NewPostgresErrorHandler()
	doaRepository = doaResultRepo.NewDOAResultRepository(db)
	economicActivityRepository = economicalActivityRepo.NewEconomicActivityRepo(db)
	offerRepository = offerRepo.NewOfferRepository(db, dbErrorHandler)
	overrideRepository = overrideRepo.NewOverrideSqlRepository(db)
	profileRepository = profileRepo.New(db, dbErrorHandler)
	personRepository = personRepo.New(db, dbErrorHandler)
	stateRepository = stateRepo.NewStateSqlRepository(db, dbErrorHandler)
	contractRepository = contractRepo.New(db, dbErrorHandler)
}

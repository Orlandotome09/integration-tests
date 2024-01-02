package _init

import (
	"context"

	limitEventTopicPublisher "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit"
	limitMessageTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit/translator"
	treeAdapterEventTopicPublisher "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter"
	treeAdapterMessageTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/accountTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/addressTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/companyTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/individualTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/timeTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account"
	accountClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address"
	addressClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http"
	registrationAddressesTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner"
	partnerClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http"
	partnerTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/engines/posProcessor"
	"bitbucket.org/bexstech/temis-compliance/src/infra/queues/pubsubPublisher"
)

func buildPosProcessor() interfaces.PosProcessor {
	ctx := context.Background()

	treeAdapterPublisher := pubsubPublisher.New(ctx, bexsDigitalPubsubClient, treeAdapterTopic)
	limitPublisher := pubsubPublisher.New(ctx, bexsDigitalPubsubClient, limitTopic)
	limitEventPublisher := limitEventTopicPublisher.New(limitPublisher, limitMessageTranslator.New())

	accountClientInstance := accountClient.New(webClient, temisRegistrationHost)
	addressClientInstance := addressClient.New(webClient, temisRegistrationHost)
	partnerClientInstance := partnerClient.New(webClient, temisRegistrationHost)

	partnerAdapter := partner.NewPartnerAdapter(partnerClientInstance, partnerTranslator.New())
	accountService := account.NewAccountAdapter(accountClientInstance)
	addressService := address.NewAddressAdapter(addressClientInstance, registrationAddressesTranslator.New())

	treeAdapterTranslator := treeAdapterMessageTranslator.New(timeTranslator.NewTimeTranslator(),
		accountTranslator.NewAccountTranslator(),
		addressTranslator.NewAddressTranslator(),
		individualTranslator.NewIndividualTranslator(timeTranslator.NewTimeTranslator()),
		companyTranslator.NewCompanyTranslator(timeTranslator.NewTimeTranslator()))

	treeAdapterEventPublisher := treeAdapterEventTopicPublisher.New(treeAdapterPublisher, treeAdapterTranslator)

	return posProcessor.NewPosProcessor(partnerAdapter, accountService, addressService, treeAdapterEventPublisher, limitEventPublisher)
}

package _init

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document"
	documentClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http"
	documentTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/translator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile"
	documentFileClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http"
	documentFileTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/contract/factory"
)

func buildContractRulesFactory() interfaces.ContractRulesFactory {
	documentClientInstance := documentClient.New(webClient, temisRegistrationHost)
	documentFileClientInstance := documentFileClient.New(webClient, temisRegistrationHost)

	documentService := document.NewDocumentAdapter(documentClientInstance, documentTranslator.New())
	documentFileService := documentFile.NewDocumentFileAdapter(documentFileClientInstance, documentFileTranslator.New())

	return contractrulesfactory.New(documentService, documentFileService, buildStateService())
}

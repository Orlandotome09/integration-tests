package blacklistAdapter

import (
	blacklistClient "bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http"
	blacklistTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/translator"
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type blacklistAdapter struct {
	blacklistClient blacklistClient.BlackListHttpClient
	translator      blacklistTranslator.BlacklistTranslator
}

func New(blacklistClient blacklistClient.BlackListHttpClient,
	translator blacklistTranslator.BlacklistTranslator) _interfaces.BlacklistAdapter {
	return &blacklistAdapter{
		blacklistClient: blacklistClient,
		translator:      translator,
	}
}

func (ref *blacklistAdapter) Search(documentNumber, partnerId string) (*entity.BlacklistStatus, bool, error) {
	response, exists, err := ref.blacklistClient.Search(documentNumber, partnerId)

	blacklistStatus := ref.translator.ToDomain(&response)

	return blacklistStatus, exists, err
}

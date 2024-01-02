package blacklistAdapter

import (
	blacklistClient "bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http/dto"
	blacklistTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/translator"
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"os"
	"reflect"
	"testing"
)

var (
	backlistHttpClient   *blacklistClient.MockBlacklistClient
	blackslisttranslator *blacklistTranslator.MockBlacklistTranslator
	service              _interfaces.BlacklistAdapter
)

func TestMain(m *testing.M) {
	backlistHttpClient = &blacklistClient.MockBlacklistClient{}
	blackslisttranslator = &blacklistTranslator.MockBlacklistTranslator{}
	service = New(backlistHttpClient, blackslisttranslator)
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	documentNumber := "111"
	partnerID := "222"

	response := dto.BlacklistResponse{}
	blacklistStatus := &entity.BlacklistStatus{}

	backlistHttpClient.On("Search", documentNumber, partnerID).Return(response, true, nil)
	blackslisttranslator.On("ToDomain", &response).Return(blacklistStatus)

	received, exists, err := service.Search(documentNumber, partnerID)

	if !reflect.DeepEqual(blacklistStatus, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", blacklistStatus, received)
	}

	if !reflect.DeepEqual(true, exists) {
		t.Errorf("\nExpected: %v \nGot: %v\n", true, exists)
	}

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", nil, err)
	}
}

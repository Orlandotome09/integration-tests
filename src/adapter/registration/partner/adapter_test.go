package partner

import (
	partnerClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http/contracts"
	partnerTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"os"
	"reflect"
	"testing"
)

var (
	client     *partnerClient.MockPartnerClient
	translator *partnerTranslator.MockPartnerTranslator
	service    interfaces.PartnerAdapter
)

func TestMain(m *testing.M) {
	client = &partnerClient.MockPartnerClient{}
	translator = &partnerTranslator.MockPartnerTranslator{}
	service = NewPartnerAdapter(client, translator)
	os.Exit(m.Run())
}

func TestGetActive(t *testing.T) {
	partnerID := "11111"

	response := &contracts.PartnerResponse{PartnerID: "2222", Status: PartnerActive}
	partner := entity.Partner{}

	client.On("Get", partnerID).Return(response, nil)
	translator.On("Translate", *response).Return(partner)

	expected := &partner
	received, err := service.GetActive(partnerID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGetActive_NoContent(t *testing.T) {
	partnerID := "22222"

	var response *contracts.PartnerResponse = nil

	client.On("Get", partnerID).Return(response, nil)

	var expected *entity.Partner = nil
	received, err := service.GetActive(partnerID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err == nil {
		t.Errorf("\nExpected err not nil")
	}
}

func TestGetActive_NotActive(t *testing.T) {
	partnerID := "33333"

	response := &contracts.PartnerResponse{PartnerID: "2222", Status: "aaa"}

	client.On("Get", partnerID).Return(response, nil)

	var expected *entity.Partner = nil
	received, err := service.GetActive(partnerID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err == nil {
		t.Errorf("\nExpected err not nil")
	}
}

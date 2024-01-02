package partnerTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := New()

	response := contracts.PartnerResponse{
		PartnerID:      "1111",
		DocumentNumber: "2222",
		Name:           "José",
		Status:         "ACTIVE",
		LogoImageUrl:   "aaaa.xxxx/cccc",
		Config: contracts.Config{
			CustomerSegregationType: "BY_PARTNER",
			UseCallbackV2:           true,
		},
	}

	expected := entity.Partner{
		PartnerID:      "1111",
		DocumentNumber: "2222",
		Name:           "José",
		Status:         "ACTIVE",
		LogoImageUrl:   "aaaa.xxxx/cccc",
		Config: &entity.PartnerConfig{
			CustomerSegregationType: "BY_PARTNER",
			UseCallbackV2:           &response.Config.UseCallbackV2,
		},
	}

	received := translator.Translate(response)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateStatus(t *testing.T) {
	status := StatusActive

	expected := values2.PartnerStatusActive
	received := translateStatus(status)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateStatus_NotActive(t *testing.T) {
	status := "NOTHING"

	var expected values2.PartnerStatus = ""
	received := translateStatus(status)
	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateSegrationType(t *testing.T) {
	segregationType := SegregationByPartner

	expected := values2.SegregationTypeByPartner
	received := translateSegregationType(segregationType)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateSegrationType_ByMerchant(t *testing.T) {
	segrationType := SegregationByMerchant

	expected := values2.SegregationTypeByMerchant
	received := translateSegregationType(segrationType)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateSegrationType_NoSegregationType(t *testing.T) {
	segrationType := "NOTHING"

	var expected values2.SegregationType = ""
	received := translateSegregationType(segrationType)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

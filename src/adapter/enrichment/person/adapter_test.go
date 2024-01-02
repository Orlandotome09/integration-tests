package personAdapter

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"

	"bitbucket.org/bexstech/temis-compliance/src/adapter"
	enrichedIndividualContracts "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/individual/contracts"
	enrichedLegalEntityContracts "bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/legalEntity/contracts"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
)

var (
	httpClient *adapter.MockHttpClient
	service    interfaces.BureauService
	headers    map[string]string
)

func TestMain(m *testing.M) {
	httpClient = &adapter.MockHttpClient{}
	service = NewPersonAdapter(httpClient)
	headers = map[string]string{"Offer-Type": "PAY_OFFER", "Partner-Id": "PAY"}
	os.Exit(m.Run())
}

func TestGetBureauStatus_ForIndividual(t *testing.T) {
	documentNumber := "96882069034"
	person := entity.Person{
		DocumentNumber: documentNumber,
		PersonType:     values.PersonTypeIndividual,
		PartnerID:      "PAY",
		OfferType:      "PAY_OFFER",
		ProfileID:      uuid.New(),
		EntityID:       uuid.New(),
	}

	response := &enrichedIndividualContracts.IndividualResponse{Situation: 1}
	status := &entity.EnrichedInformation{
		BureauStatus: "REGULAR",
	}

	individualBytes := new(bytes.Buffer)
	json.NewEncoder(individualBytes).Encode(response)

	httpClient.On("Get", GetIndividualPath+documentNumber, "", headers).Return(individualBytes.Bytes(), nil)

	expected := status
	received, err := service.GetBureauStatus(person)

	if err != nil {
		t.Errorf("\nExpected error nil \n")
	}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestGetBureauStatus_ForIndividual_notFoundInBureau(t *testing.T) {
	documentNumber := "25031902868"
	person := entity.Person{
		DocumentNumber: documentNumber,
		PersonType:     values.PersonTypeIndividual,
		PartnerID:      "PAY",
		OfferType:      "PAY_OFFER",
		ProfileID:      uuid.New(),
		EntityID:       uuid.New(),
	}

	httpClient.On("Get", GetIndividualPath+documentNumber, "", headers).Return(nil, nil)

	var expected *entity.EnrichedInformation = nil
	received, err := service.GetBureauStatus(person)

	if err != nil {
		t.Errorf("\nExpected error nil \n")
	}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestGetBureauStatus_ForLegalEntity(t *testing.T) {
	legalEntityID := "09684205000108"
	person := entity.Person{
		DocumentNumber: legalEntityID,
		PersonType:     values.PersonTypeCompany,
		PartnerID:      "PAY",
		OfferType:      "PAY_OFFER",
		ProfileID:      uuid.New(),
		EntityID:       uuid.New(),
	}

	legalEntityResponse := &enrichedLegalEntityContracts.LegalEntityResponse{Situation: 2}
	status := &entity.EnrichedInformation{BureauStatus: "REGULAR", EnrichedCompany: entity.EnrichedCompany{BoardOfDirectors: []entity.Director{}}}

	legalEntityBytes := new(bytes.Buffer)
	json.NewEncoder(legalEntityBytes).Encode(legalEntityResponse)

	httpClient.On("Get", GetLegalEntityPath+legalEntityID, "", headers).Return(legalEntityBytes.Bytes(), nil)

	expected := status
	received, err := service.GetBureauStatus(person)

	if err != nil {
		t.Errorf("\nExpected error nil \n")
	}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestGetBureauStatus_ForLegalEntity_notFoundInBureau(t *testing.T) {
	legalEntityID := "11114931341829"
	person := entity.Person{
		DocumentNumber: legalEntityID,
		PersonType:     values.PersonTypeCompany,
		PartnerID:      "PAY",
		OfferType:      "PAY_OFFER",
		ProfileID:      uuid.New(),
		EntityID:       uuid.New(),
	}

	httpClient.On("Get", GetLegalEntityPath+legalEntityID, "", headers).Return(nil, nil)

	var expected *entity.EnrichedInformation = nil
	received, err := service.GetBureauStatus(person)

	if err != nil {
		t.Errorf("\nExpected error nil \n")
	}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

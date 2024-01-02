package enricherClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http/contracts"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetEnrichedPerson(t *testing.T) {
	profileID := uuid.New()
	partnerID := uuid.New()
	documentNumber := "15928485085"
	enricherRequest := enricherContracts.EnricherRequest{
		ProfileID:  profileID.String(),
		PersonType: values2.PersonTypeIndividual,
		OfferType:  "MAQUININHA_01",
		PartnerID:  partnerID.String(),
		RoleType:   values2.RoleTypeCustomer,
	}
	response := &enricherContracts.EnricherResponse{

		Person: enricherContracts.Person{
			EntityID:       profileID,
			Role:           values2.RoleTypeCustomer,
			Type:           values2.PersonTypeIndividual,
			Name:           "JOSE DA SILVA",
			DocumentNumber: documentNumber,
			Individual: &enricherContracts.IndividualResponse{
				BirthDate: "25/01/1980",
				Situation: 1,
			},
		},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	expected := response
	received, err := client.GetEnrichedPerson(enricherRequest, documentNumber)

	assert.Nil(t, err)
	assert.NotNil(t, received)
	assert.Equal(t, received, expected)
}

func TestGetEnrichedPerson_BadRequest(t *testing.T) {
	profileID := uuid.New()
	documentNumber := "15928485085"
	enricherRequest := enricherContracts.EnricherRequest{}
	response := &enricherContracts.EnricherResponse{
		Person: enricherContracts.Person{EntityID: profileID},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 400, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	received, err := client.GetEnrichedPerson(enricherRequest, documentNumber)

	assert.Nil(t, received)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Status Code: 400")
}

func TestGetEnrichedPerson_NoContent(t *testing.T) {
	profileID := uuid.New()
	partnerID := uuid.New()
	documentNumber := "15928485085"
	enricherRequest := enricherContracts.EnricherRequest{
		ProfileID:  profileID.String(),
		PersonType: values2.PersonTypeIndividual,
		OfferType:  "OFFER_NOT_EXIST",
		PartnerID:  partnerID.String(),
		RoleType:   values2.RoleTypeCustomer,
	}

	server := mockServer("GET", 204, nil, nil)
	defer server.Close()

	client := New(server.Client(), server.URL)

	received, err := client.GetEnrichedPerson(enricherRequest, documentNumber)

	assert.Nil(t, received)
	assert.Nil(t, err)
}

func TestGetEnrichedPerson_Error(t *testing.T) {
	profileID := uuid.New()
	partnerID := uuid.New()
	documentNumber := "15928485085"
	enricherRequest := enricherContracts.EnricherRequest{
		ProfileID:  profileID.String(),
		PersonType: values2.PersonTypeIndividual,
		OfferType:  "MAQUININHA_01",
		PartnerID:  partnerID.String(),
		RoleType:   values2.RoleTypeCustomer,
	}
	response := &enricherContracts.EnricherResponse{
		Person: enricherContracts.Person{EntityID: profileID},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 500, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	received, err := client.GetEnrichedPerson(enricherRequest, documentNumber)

	assert.Nil(t, received)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Status Code: 500")
}

//-------------------------------MockServer--------------------------------------
func mockServer(method string, responseStatus int, requestBody []byte, responseBody []byte) *httptest.Server {

	req := httptest.NewRequest(method, "http://test", bytes.NewBuffer(requestBody))
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responseStatus)
		w.Write(responseBody)
	})
	handler.ServeHTTP(responseRecorder, req)
	server := httptest.NewServer(handler)

	return server
}

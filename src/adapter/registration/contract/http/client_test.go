package contractClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http/contracts"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	id := uuid.New()
	response := &contracts.GetContractResponse{ContractID: &id}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	expected := response
	received, exists, err := client.Get(&id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if exists != true {
		t.Errorf("\nExpected exists true")
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGet_NotFound(t *testing.T) {
	id := uuid.New()
	response := &contracts.GetContractResponse{ContractID: &id}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 404, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	received, exists, err := client.Get(&id)

	if received != nil {
		t.Errorf("\nExpected received nil")
	}

	if exists != false {
		t.Errorf("\nExpected exists false")
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
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

package ownershipStructureClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/http/contracts"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetOwnershipStructure(t *testing.T) {
	response := &contracts.OwnershipStructureResponse{LegalEntityID: "1234"}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	expected := response
	received, err := client.GetOwnershipStructure("1111", "PAY_CATALOG", "PAY")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected no err (nil)")
	}
}

func TestGetOwnershipStructure_NoContent(t *testing.T) {
	response := &contracts.OwnershipStructureResponse{LegalEntityID: "1234"}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 204, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	var expected *contracts.OwnershipStructureResponse = nil
	received, err := client.GetOwnershipStructure("1111", "PAY_CATALOG", "PAY")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected no err (nil)")
	}
}

func TestGetOwnershipStructure_Error(t *testing.T) {
	response := &contracts.OwnershipStructureResponse{LegalEntityID: "1234"}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 500, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	var expected *contracts.OwnershipStructureResponse = nil
	received, err := client.GetOwnershipStructure("1111", "PAY_CATALOG", "PAY")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err == nil {
		t.Errorf("\nExpected err not nil")
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

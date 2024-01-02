package adapter

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGet_shouldGetSuccessfully(t *testing.T) {
	response := &entity.OwnershipStructure{FinalBeneficiariesCount: 1111}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := NewHttpClient(server.Client(), server.URL)

	expected := responseBody
	received, err := client.Get("/xxx", "?field=xxx", nil)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected no err (nil)")
	}
}

func TestGet_shouldNotFind_404(t *testing.T) {
	response := &entity.OwnershipStructure{FinalBeneficiariesCount: 2222}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 404, nil, responseBody)
	defer server.Close()

	client := NewHttpClient(server.Client(), server.URL)

	var expected []byte
	received, err := client.Get("/xxx", "?field=xxx", nil)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected no err (nil)")
	}
}

func TestGet_shouldNotFind_204(t *testing.T) {
	response := &entity.OwnershipStructure{FinalBeneficiariesCount: 3333}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 204, nil, responseBody)
	defer server.Close()

	client := NewHttpClient(server.Client(), server.URL)

	var expected []byte
	received, err := client.Get("/xxx", "?field=xxx", nil)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected no err (nil)")
	}
}

func TestGet_shouldFindError_bodyJson(t *testing.T) {
	response := &entity.OwnershipStructure{FinalBeneficiariesCount: 4444}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 500, nil, responseBody)
	defer server.Close()

	client := NewHttpClient(server.Client(), server.URL)

	var expected []byte
	received, err := client.Get("/xxx", "?field=xxx", nil)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err == nil {
		t.Errorf("\nExpected err (not nil)")
	}
}

func TestGet_shouldFindError_bodyNil(t *testing.T) {

	server := mockServer("GET", 500, nil, nil)
	defer server.Close()

	client := NewHttpClient(server.Client(), server.URL)

	var expected []byte
	received, err := client.Get("/xxx", "?field=xxx", nil)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err == nil {
		t.Errorf("\nExpected err (not nil)")
	}
}

func TestGet_shouldFindError_bodyString(t *testing.T) {
	responseBody, _ := json.Marshal("server error")

	server := mockServer("GET", 500, nil, responseBody)
	defer server.Close()

	client := NewHttpClient(server.Client(), server.URL)

	var expected []byte
	received, err := client.Get("/xxx", "?field=xxx", nil)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err == nil {
		t.Errorf("\nExpected err (not nil)")
	}
}

func TestConvertByteToString(t *testing.T) {
	response := &entity.OwnershipStructure{FinalBeneficiariesCount: 4444}
	responseBody, _ := json.Marshal(response)

	fmt.Println(string(responseBody))
}

// -------------------------------MockServer--------------------------------------
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

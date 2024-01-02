package documentClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http/contracts"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	documentID := uuid.MustParse("111d0cc3-f173-4a78-86c8-5b9de2fbf8d2")
	response := &contracts.DocumentResponse{DocumentID: documentID}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	expected := response
	received, err := client.Get("111")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGet_NoContent(t *testing.T) {
	server := mockServer("GET", 404, nil, nil)
	defer server.Close()

	client := New(server.Client(), server.URL)

	var expected *contracts.DocumentResponse = nil
	received, err := client.Get("111")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGet_ServerError(t *testing.T) {
	server := mockServer("GET", 500, nil, nil)
	defer server.Close()

	client := New(server.Client(), server.URL)

	var expected *contracts.DocumentResponse = nil
	received, err := client.Get("111")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err == nil {
		t.Errorf("\nExpected err not nil")
	}
}

func TestSearchByEntityID(t *testing.T) {
	documentID1 := uuid.MustParse("333d0cc3-f173-4a78-86c8-5b9de2fbf8d2")
	documentID2 := uuid.MustParse("444d0cc3-f173-4a78-86c8-5b9de2fbf8d2")
	response := []contracts.DocumentResponse{
		{DocumentID: documentID1}, {DocumentID: documentID2},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	expected := response
	received, err := client.SearchByEntityID("111")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestSearchByEntityID_ServerError(t *testing.T) {

	server := mockServer("GET", 500, nil, nil)
	defer server.Close()

	client := New(server.Client(), server.URL)

	var expected []contracts.DocumentResponse = nil
	received, err := client.SearchByEntityID("111")

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

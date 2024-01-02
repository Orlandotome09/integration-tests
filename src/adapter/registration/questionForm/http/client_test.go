package questionFormClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/questionForm/http/contracts"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	response := &contracts.QuestionFormResponse{QuestionFormID: "1111"}
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

	var expected *contracts.QuestionFormResponse = nil
	received, err := client.Get("111")

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
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

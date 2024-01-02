package restrictiveListsHttpClient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSearchInInternalList(t *testing.T) {
	response := InternalListResponse{
		{
			DocumentNumber: uuid.New().String(),
		},
		{
			DocumentNumber: uuid.New().String(),
		},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := NewRestrictiveListHttpClient(server.Client(), server.URL)

	expected := response
	received, err := client.SearchInternalList("111", "name")

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestSearchInInternalList_NoMatches(t *testing.T) {
	response := InternalListResponse{}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := NewRestrictiveListHttpClient(server.Client(), server.URL)

	expected := response
	received, err := client.SearchInternalList("111", "name")

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestSearchInInternalList_ServerError(t *testing.T) {

	server := mockServer("GET", 500, nil, nil)
	defer server.Close()

	client := NewRestrictiveListHttpClient(server.Client(), server.URL)

	received, err := client.SearchInternalList("111", "name")

	assert.NotNil(t, err)
	assert.Nil(t, received)
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

func TestSearchPepList(t *testing.T) {
	response := PepResponse{
		DocumentNumber: "123",
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := NewRestrictiveListHttpClient(server.Client(), server.URL)

	expected := &response
	received, err := client.SearchPepList("123")

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestSearchPepList_NotFound(t *testing.T) {
	response := PepResponse{
		DocumentNumber: "444",
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 204, nil, responseBody)
	defer server.Close()

	client := NewRestrictiveListHttpClient(server.Client(), server.URL)

	received, err := client.SearchPepList("123")

	assert.Nil(t, err)
	assert.Nil(t, received)
}

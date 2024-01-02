package notificationRecipientClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearch(t *testing.T) {
	response := []contracts.NotificationRecipientResponse{
		{NotificationType: values.NotificationTypePostWarnings.ToString()},
		{NotificationType: values.NotificationTypeSentOP.ToString()},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL)

	expected := response
	received, err := client.Search("111")

	assert.Equal(t, expected, received)
	assert.Nil(t, err)
}

func TestSearch_ServerError(t *testing.T) {
	server := mockServer("GET", 500, nil, nil)
	defer server.Close()

	client := New(server.Client(), server.URL)

	received, err := client.Search("111")

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(received))
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

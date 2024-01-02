package watchlistClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/auth"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http/dto"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	documentNumber := "222"
	firstName := "Aaa"
	lastName := "Bbbb"
	fullName := "Aaa Bbbb"
	birthYear := 2000

	response := []dto.WatchlistResponse{
		{Title: "Title1"}, {Title: "Title2"}, {Title: "Title3"},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL, nil)

	expected := response
	received, err := client.SearchIndividual(documentNumber, firstName, lastName, fullName, "BR", birthYear)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil, got: %v", err)
	}
}

func TestSearch_NotFound(t *testing.T) {
	documentNumber := "222"
	firstName := "Aaa"
	lastName := "Bbbb"
	fullName := "Aaa Bbbb"
	birthYear := 2000

	response := []dto.WatchlistResponse{
		{Title: "Title1"}, {Title: "Title2"}, {Title: "Title3"},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 204, nil, responseBody)
	defer server.Close()

	client := New(server.Client(), server.URL, nil)

	received, err := client.SearchIndividual(documentNumber, firstName, lastName, fullName, "BR", birthYear)

	if received != nil {
		t.Errorf("\nExpected: %v \nGot: %v\n", nil, received)
	}

	if err == nil {
		t.Errorf("\nExpected err not nil")
	}
}

func TestSearchWithAuth(t *testing.T) {
	authRepo := &auth.MockAuthRepository{}

	documentNumber := "222"
	firstName := "Aaa"
	lastName := "Bbbb"
	fullName := "Aaa Bbbb"
	birthYear := 2000

	response := []dto.WatchlistResponse{
		{Title: "Title1"}, {Title: "Title2"}, {Title: "Title3"},
	}
	responseBody, _ := json.Marshal(response)

	server := mockServer("GET", 200, nil, responseBody)
	defer server.Close()

	token := "token"
	authRepo.On("GetAccessToken").Return(token, nil)

	client := New(server.Client(), server.URL, authRepo)

	expected := response
	received, err := client.SearchIndividual(documentNumber, firstName, lastName, fullName, "BR", birthYear)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil, got: %v", err)
	}

	mock.AssertExpectationsForObjects(t, authRepo)
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

package http

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCadastralValidationConfig(t *testing.T) {
	ctx := context.TODO()
	t.Run("should get response", func(t *testing.T) {
		catalogsCache := &mocks.Cache{}

		personType := uuid.NewString()
		roleType := uuid.NewString()
		offerType := uuid.NewString()
		cacheKey := personType + roleType + offerType
		response := CadastralValidationConfigResponse{
			{
				PersonType: uuid.NewString(),
				RoleType:   uuid.NewString(),
				OfferType:  uuid.NewString(),
				PartnerID:  uuid.NewString(),
			},
		}
		responseBody, _ := json.Marshal(response)

		server := mockServer("GET", 200, nil, responseBody)
		defer server.Close()

		client := NewTemisConfigHttpClient(server.Client(), server.URL, catalogsCache)

		catalogsCache.On("Get", cacheKey).Return(nil)
		catalogsCache.On("Save", cacheKey, response, cacheTTL).Run(nil)

		expected := response
		received, err := client.GetCadastralValidationConfig(ctx, personType, roleType, offerType)

		assert.Nil(t, err)
		assert.Equal(t, expected, received)
		mock.AssertExpectationsForObjects(t, catalogsCache)
	})

	t.Run("should get error 500", func(t *testing.T) {
		catalogsCache := &mocks.Cache{}

		personType := uuid.NewString()
		roleType := uuid.NewString()
		offerType := uuid.NewString()
		cacheKey := personType + roleType + offerType
		response := CadastralValidationConfigResponse{
			{
				PersonType: uuid.NewString(),
				RoleType:   uuid.NewString(),
				OfferType:  uuid.NewString(),
				PartnerID:  uuid.NewString(),
			},
		}
		responseBody, _ := json.Marshal(response)

		server := mockServer("GET", 500, nil, responseBody)
		defer server.Close()

		catalogsCache.On("Get", cacheKey).Return(nil)

		client := NewTemisConfigHttpClient(server.Client(), server.URL, catalogsCache)

		var expected CadastralValidationConfigResponse
		received, err := client.GetCadastralValidationConfig(ctx, personType, roleType, offerType)

		assert.ErrorContains(t, err, "Error requesting cadastral validation config.")
		assert.Equal(t, expected, received)
		mock.AssertExpectationsForObjects(t, catalogsCache)
		catalogsCache.AssertNumberOfCalls(t, "Save", 0)
	})
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

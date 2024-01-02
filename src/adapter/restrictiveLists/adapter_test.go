package restrictivelists

import (
	"fmt"
	"testing"
	"time"

	restrictiveListsHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/restrictiveLists/http"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOccurrenceInBlackList(t *testing.T) {
	httpClient := &restrictiveListsHttpClient.MockRestrictiveListsHttpClient{}
	adapter := NewRestrictiveListsAdapter(httpClient)

	document := uuid.New().String()
	name := "name"
	createdAt := time.Now()
	internalList := restrictiveListsHttpClient.InternalListResponse{
		{
			Author:        "John Doe",
			Justification: "internal list",
			CreatedAt:     createdAt,
		},
	}

	httpClient.On("SearchInternalList", document, name).Return(internalList, nil)

	received, err := adapter.OccurrenceInBlackList(document, name)

	expected := &entity.BlacklistStatus{
		Justification: entity.Justification{
			AddedAt:  internalList[0].CreatedAt,
			Author:   internalList[0].Author,
			Comments: []string{internalList[0].Justification},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestOccurrenceInBlackList_NoOccurences(t *testing.T) {
	httpClient := &restrictiveListsHttpClient.MockRestrictiveListsHttpClient{}
	adapter := NewRestrictiveListsAdapter(httpClient)

	document := uuid.New().String()
	name := "name"
	internalList := restrictiveListsHttpClient.InternalListResponse{}

	httpClient.On("SearchInternalList", document, name).Return(internalList, nil)

	received, err := adapter.OccurrenceInBlackList(document, name)

	assert.Nil(t, err)
	assert.Nil(t, received)
}

func TestOccurrenceInBlackList_ErrSearchingInternalList(t *testing.T) {
	httpClient := &restrictiveListsHttpClient.MockRestrictiveListsHttpClient{}
	adapter := NewRestrictiveListsAdapter(httpClient)

	document := uuid.New().String()
	name := "name"
	var internalList restrictiveListsHttpClient.InternalListResponse
	errorSearching := fmt.Errorf("error searching")

	httpClient.On("SearchInternalList", document, name).Return(internalList, errorSearching)

	received, err := adapter.OccurrenceInBlackList(document, name)

	assert.Nil(t, received)
	assert.Equal(t, err.Error(), errorSearching.Error())
}

func TestOccurrenceInPepList(t *testing.T) {
	httpClient := &restrictiveListsHttpClient.MockRestrictiveListsHttpClient{}
	adapter := NewRestrictiveListsAdapter(httpClient)

	documentNumber := "123"
	pepResponse := restrictiveListsHttpClient.PepResponse{
		DocumentNumber: documentNumber,
	}

	httpClient.On("SearchPepList", documentNumber).Return(&pepResponse, nil)

	received, err := adapter.OccurrenceInPepList(documentNumber)

	expected := &entity.PepInformation{DocumentNumber: documentNumber}

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestOccurrenceInPepList_NotFound(t *testing.T) {
	httpClient := &restrictiveListsHttpClient.MockRestrictiveListsHttpClient{}
	adapter := NewRestrictiveListsAdapter(httpClient)

	documentNumber := "123"
	var pepResponse *restrictiveListsHttpClient.PepResponse = nil

	httpClient.On("SearchPepList", documentNumber).Return(pepResponse, nil)

	received, err := adapter.OccurrenceInPepList(documentNumber)

	var expected *entity.PepInformation = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

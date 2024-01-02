package doaAdapter

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/doa/contracts"
	translator "bitbucket.org/bexstech/temis-compliance/src/adapter/doa/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	doaMockHttpClient *adapter.MockHttpClient
	doaMockTranslator *translator.MockDOATranslator
	doaMockAdapter    interfaces.DOAAdapter
)

func TestMain(m *testing.M) {
	doaMockHttpClient = &adapter.MockHttpClient{}
	doaMockTranslator = &translator.MockDOATranslator{}
	doaMockAdapter = NewDOAAdapter(doaMockHttpClient, doaMockTranslator)
	os.Exit(m.Run())
}

func TestRequestExtraction(t *testing.T) {

	frontFile := &entity.DocumentFile{FileID: uuid.New()}
	frontFileURI := "frontFileURI"
	backFile := &entity.DocumentFile{FileID: uuid.New()}
	backFileURI := "backFileURL"

	time := time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)
	doc := &entity.Document{
		DocumentSubType: string(values.DocumentSubTypeRg),
		DocumentFields:  entity.DocumentFields{Number: "12345", IssueDate: time.String(), Name: "Maria"},
	}

	profileID := uuid.New()

	doaExtractionRequest := &contracts.DOAExtractionRequest{}

	response := &contracts.DOAExtractionResponse{
		Message:   "message test",
		RequestID: uuid.New(),
	}
	doaBytes := new(bytes.Buffer)
	json.NewEncoder(doaBytes).Encode(response)

	doaExtraction := &entity.DOAExtraction{
		Message:   response.Message,
		RequestID: response.RequestID,
	}

	doaMockTranslator.On("FromDomain", frontFile, frontFileURI, backFile, backFileURI, doc, profileID).Return(doaExtractionRequest, nil)
	doaMockHttpClient.On("Post", GetRequestExtractionPath, doaExtractionRequest).Return(doaBytes.Bytes(), nil)
	doaMockTranslator.On("ToDomain", doaBytes.Bytes()).Return(doaExtraction, nil)

	expected := doaExtraction
	received, err := doaMockAdapter.RequestExtraction(frontFile, frontFileURI, backFile, backFileURI, doc, profileID)

	if err != nil {
		t.Errorf("\nExpected error nil \n")
	}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

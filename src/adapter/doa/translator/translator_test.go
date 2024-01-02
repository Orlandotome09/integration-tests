package doaTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/doa/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestFromDomain(t *testing.T) {
	temisComplianceHost := "temisComplianceHost"
	translator := New(temisComplianceHost)

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

	file1 := contracts.FileParams{
		FileID:   frontFile.FileID.String(),
		FileSide: contracts.ToFileSide(frontFile.FileSide),
		FileURI:  frontFileURI,
		FileName: frontFile.FileID.String(),
	}

	file2 := contracts.FileParams{
		FileID:   backFile.FileID.String(),
		FileSide: contracts.ToFileSide(backFile.FileSide),
		FileURI:  backFileURI,
		FileName: backFile.FileID.String(),
	}

	expected := &contracts.DOAExtractionRequest{
		ProfileID:      profileID,
		Type:           contracts.ToDocumentType(doc.DocumentSubType),
		Metadata:       contracts.ToMetadata(doc.DocumentFields),
		FileParams:     []contracts.FileParams{file1, file2},
		CallbackParams: &contracts.CallbackParams{URL: temisComplianceHost + "/doa/callback/"},
	}

	received, err := translator.FromDomain(frontFile, frontFileURI, backFile, backFileURI, doc, profileID)

	assert.Nil(t, err)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestToDomain(t *testing.T) {
	translator := New("")

	response := &contracts.DOAExtractionResponse{
		Message:   "message test",
		RequestID: uuid.New(),
	}

	expected := &entity.DOAExtraction{
		Message:   response.Message,
		RequestID: response.RequestID,
	}

	doaBytes := new(bytes.Buffer)
	json.NewEncoder(doaBytes).Encode(response)

	received, err := translator.ToDomain(doaBytes.Bytes())

	assert.Nil(t, err)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

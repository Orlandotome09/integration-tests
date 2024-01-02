package document

import (
	documentClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http/contracts"
	documentTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
)

var (
	client     *documentClient.MockDocumentClient
	translator *documentTranslator.MockDocumentTranslator
	adapter    interfaces.DocumentAdapter
)

func TestMain(m *testing.M) {
	client = &documentClient.MockDocumentClient{}
	translator = &documentTranslator.MockDocumentTranslator{}
	adapter = NewDocumentAdapter(client, translator)
	os.Exit(m.Run())
}

func TestGetByID(t *testing.T) {
	id := "1111"

	response := &contracts.DocumentResponse{}
	document := entity.Document{}

	client.On("Get", id).Return(response, nil)
	translator.On("Translate", *response).Return(document, nil)

	expected := &document
	received, err := adapter.GetByID(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGetByID_NoContent(t *testing.T) {
	id := "4444"

	var response *contracts.DocumentResponse = nil

	client.On("Get", id).Return(response, nil)

	var expected *entity.Document = nil
	received, err := adapter.GetByID(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestFindByEntityIDAndDocumentType(t *testing.T) {
	id := "2222"
	documentType := "ID"

	documentID1 := uuid.MustParse("111ca7c6-d667-4a98-8a9d-d5cae928cd7f")
	documentID2 := uuid.MustParse("222ca7c6-d667-4a98-8a9d-d5cae928cd7f")
	documentID3 := uuid.MustParse("333ca7c6-d667-4a98-8a9d-d5cae928cd7f")
	responses := []contracts.DocumentResponse{}
	documents := []entity.Document{
		{DocumentID: documentID1, DocumentType: "ID"},
		{DocumentID: documentID2, DocumentType: "XX"},
		{DocumentID: documentID3, DocumentType: "ID"},
	}

	client.On("SearchByEntityID", id).Return(responses, nil)

	translator.On("TranslateAll", responses).Return(documents, nil)

	expected := []entity.Document{documents[0], documents[2]}
	received, err := adapter.FindByEntityIDAndDocumentType(id, documentType)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestFilterByDocumentType(t *testing.T) {
	documentType := "ID"

	documentID1 := uuid.MustParse("111ca7c6-d667-4a98-8a9d-d5cae928cd22")
	documentID2 := uuid.MustParse("222ca7c6-d667-4a98-8a9d-d5cae928cd22")
	documentID3 := uuid.MustParse("333ca7c6-d667-4a98-8a9d-d5cae928cd22")

	documents := []entity.Document{
		{DocumentID: documentID1, DocumentType: "ID"},
		{DocumentID: documentID2, DocumentType: "XX"},
		{DocumentID: documentID3, DocumentType: "ID"},
	}

	expected := []entity.Document{documents[0], documents[2]}
	received := filterByDocumentType(documentType, documents)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

}

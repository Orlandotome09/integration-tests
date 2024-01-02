package documentTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := New()
	documentID := uuid.MustParse("111ca7c6-d667-4a98-8a9d-d5cae928cd7f")
	profileID := uuid.MustParse("2222d7cd-ec6d-450d-b6e3-451eab3e62a4")
	response := contracts.DocumentResponse{
		DocumentID: documentID,
		EntityID:   profileID,
		Type:       "IDENTIFICATION",
		SubType:    "RG",
		DocumentFields: contracts.DocumentFields{
			Name:      "Pedro",
			Number:    "12345",
			IssueDate: "121212",
		},
		ExpirationDate: "2024-12-30",
	}

	expected := entity.Document{
		DocumentID:      documentID,
		EntityID:        profileID,
		DocumentType:    "IDENTIFICATION",
		DocumentSubType: "RG",
		DocumentFields: entity.DocumentFields{
			Name:      "Pedro",
			Number:    "12345",
			IssueDate: "121212",
		},
		ExpirationDate: "2024-12-30",
	}

	received := translator.Translate(response)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslate_EmptyObject(t *testing.T) {
	translator := New()

	response := contracts.DocumentResponse{}

	expected := entity.Document{}

	received := translator.Translate(response)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateAll(t *testing.T) {
	translator := New()
	documentID1 := uuid.MustParse("111ca7c6-d667-4a98-8a9d-d5cae928cd7f")
	documentID2 := uuid.MustParse("2222d7cd-ec6d-450d-b6e3-451eab3e62a4")
	profileID1 := uuid.MustParse("3333a7c6-d667-4a98-8a9d-d5cae928cd7f")
	profileID2 := uuid.MustParse("4442d7cd-ec6d-450d-b6e3-451eab3e62a4")

	responses := []contracts.DocumentResponse{
		{
			DocumentID: documentID1,
			EntityID:   profileID1,
			Type:       "IDENTIFICATION",
			SubType:    "RG",
			DocumentFields: contracts.DocumentFields{
				Name:      "AAAA",
				IssueDate: "1111",
				Number:    "1111",
			},
		},
		{
			DocumentID: documentID2,
			EntityID:   profileID2,
			Type:       "IDENTIFICATION",
			SubType:    "CNH",
			DocumentFields: contracts.DocumentFields{
				Name:      "BBBB",
				IssueDate: "3333",
				Number:    "3333",
			},
		},
	}

	expected := []entity.Document{
		{
			DocumentID:      documentID1,
			EntityID:        profileID1,
			DocumentType:    "IDENTIFICATION",
			DocumentSubType: "RG",
			DocumentFields: entity.DocumentFields{
				Name:      "AAAA",
				IssueDate: "1111",
				Number:    "1111",
			},
		},
		{
			DocumentID:      documentID2,
			EntityID:        profileID2,
			DocumentType:    "IDENTIFICATION",
			DocumentSubType: "CNH",
			DocumentFields: entity.DocumentFields{
				Name:      "BBBB",
				IssueDate: "3333",
				Number:    "3333",
			},
		},
	}

	received := translator.TranslateAll(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateAll_Empty(t *testing.T) {
	translator := New()

	responses := []contracts.DocumentResponse{{}, {}}

	expected := []entity.Document{{}, {}}

	received := translator.TranslateAll(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

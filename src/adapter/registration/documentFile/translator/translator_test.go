package documentFileTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {
	translator := New()

	time := time.Date(2000, 10, 10, 10, 10, 10, 10, time.UTC)

	response := contracts.DocumentFileResponse{
		DocumentFileID: uuid.MustParse("462ca7c6-d667-4a98-8a9d-d5cae928cd7f"),
		DocumentID:     uuid.MustParse("6172d7cd-ec6d-450d-b6e3-451eab3e62a4"),
		FileID:         uuid.MustParse("5ba8db6a-612f-4c15-9e33-e97186706a6b"),
		FileSide:       "FRENTE",
		CreatedAt:      time,
	}

	expected := entity.DocumentFile{
		DocumentFileID: &response.DocumentFileID,
		DocumentID:     response.DocumentID,
		FileID:         response.FileID,
		FileSide:       values.FileSide(response.FileSide),
		CreatedAt:      response.CreatedAt,
	}

	received := translator.Translate(response)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateAll(t *testing.T) {
	translator := New()

	responses := []contracts.DocumentFileResponse{
		{
			DocumentFileID: uuid.MustParse("462ca7c6-d667-4a98-8a9d-d5cae928cd7f"),
			DocumentID:     uuid.MustParse("5172d7cd-ec6d-450d-b6e3-451eab3e62a4"),
			FileID:         uuid.MustParse("6ba8db6a-612f-4c15-9e33-e97186706a6b"),
			FileSide:       "FRONT",
			CreatedAt:      time.Date(2000, 10, 10, 10, 10, 10, 10, time.UTC),
		},
		{
			DocumentFileID: uuid.MustParse("762ca7c6-d667-4a98-8a9d-d5cae928cd7f"),
			DocumentID:     uuid.MustParse("8172d7cd-ec6d-450d-b6e3-451eab3e62a4"),
			FileID:         uuid.MustParse("9ba8db6a-612f-4c15-9e33-e97186706a6b"),
			FileSide:       "BACK",
			CreatedAt:      time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC),
		},
	}

	expected := []entity.DocumentFile{
		{
			DocumentFileID: &responses[0].DocumentFileID,
			DocumentID:     responses[0].DocumentID,
			FileID:         responses[0].FileID,
			FileSide:       values.FileSide(responses[0].FileSide),
			CreatedAt:      responses[0].CreatedAt,
		},
		{
			DocumentFileID: &responses[1].DocumentFileID,
			DocumentID:     responses[1].DocumentID,
			FileID:         responses[1].FileID,
			FileSide:       values.FileSide(responses[1].FileSide),
			CreatedAt:      responses[1].CreatedAt,
		},
	}

	received := translator.TranslateAll(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

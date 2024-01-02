package documentFile

import (
	documentFileClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http/contracts"
	documentFileTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	documentfileClient     *documentFileClient.MockDocumentFileClient
	documentfileTranslator *documentFileTranslator.MockDocumentFileTranslator
	adapter                interfaces.DocumentFileAdapter
)

func TestMain(m *testing.M) {
	documentfileClient = &documentFileClient.MockDocumentFileClient{}
	documentfileTranslator = &documentFileTranslator.MockDocumentFileTranslator{}
	adapter = NewDocumentFileAdapter(documentfileClient, documentfileTranslator)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	documentFileID := uuid.MustParse("111a6ce0-9a19-43b8-9776-44ddab16d5b4")
	response := &contracts.DocumentFileResponse{}
	documentFile := entity.DocumentFile{}

	documentfileClient.On("Get", documentFileID.String()).Return(response, nil)
	documentfileTranslator.On("Translate", *response).Return(documentFile)

	expected := &documentFile
	received, err := adapter.Get(documentFileID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGet_NotContent(t *testing.T) {
	documentFileID := uuid.MustParse("222a6ce0-9a19-43b8-9776-44ddab16d5b4")
	var response *contracts.DocumentFileResponse = nil

	documentfileClient.On("Get", documentFileID.String()).Return(response, nil)

	var expected *entity.DocumentFile = nil
	received, err := adapter.Get(documentFileID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestFindByDocumentID(t *testing.T) {
	documentID := uuid.MustParse("321a6ce0-9a19-43b8-9776-44dddd16d5b4")

	responses := []contracts.DocumentFileResponse{}
	documentFiles := []entity.DocumentFile{}

	documentfileClient.On("SearchByDocumentId", documentID.String()).Return(responses, nil)
	documentfileTranslator.On("TranslateAll", responses).Return(documentFiles)

	expected := documentFiles
	received, err := adapter.FindByDocumentID(documentID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGetLastTwoFilesOfDocument(t *testing.T) {
	documentID := uuid.MustParse("321a6ce0-9a19-43b8-9776-55dddd16d5b4")

	documentFileID1 := uuid.MustParse("555a6ce0-9a19-43b8-9776-55ddab16d5b4")
	documentFileID2 := uuid.MustParse("666912b9-5ce1-495e-8681-55818744954a")

	time1 := time.Date(2030, 10, 10, 10, 10, 10, 10, time.UTC)
	time2 := time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)

	responses := []contracts.DocumentFileResponse{
		{DocumentFileID: documentFileID1, FileSide: frontFile, CreatedAt: time1},
		{DocumentFileID: documentFileID2, FileSide: backFile, CreatedAt: time2},
	}
	frontFileResponses := []contracts.DocumentFileResponse{responses[0]}
	backFileResponses := []contracts.DocumentFileResponse{responses[1]}

	frontFiles := []entity.DocumentFile{
		{DocumentFileID: &documentFileID1, FileSide: frontFile, CreatedAt: time1},
	}
	backFiles := []entity.DocumentFile{
		{DocumentFileID: &documentFileID2, FileSide: backFile, CreatedAt: time2},
	}

	documentfileClient.On("SearchByDocumentId", documentID.String()).Return(responses, nil)
	documentfileTranslator.On("TranslateAll", frontFileResponses).Return(frontFiles)
	documentfileTranslator.On("TranslateAll", backFileResponses).Return(backFiles)

	expectedFront := &frontFiles[0]
	expectedBack := &backFiles[0]
	receivedFront, receivedBack, err := adapter.GetLastTwoFilesOfDocument(documentID)

	if !reflect.DeepEqual(expectedFront, receivedFront) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expectedFront, receivedFront)
	}

	if !reflect.DeepEqual(expectedBack, receivedBack) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expectedBack, receivedBack)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestSplitResponsesByFileSide(t *testing.T) {
	documentFileID1 := uuid.MustParse("111a6ce0-9a19-43b8-9776-44ddab16d5b4")
	documentFileID2 := uuid.MustParse("222912b9-5ce1-495e-8681-be818744954a")
	documentFileID3 := uuid.MustParse("33387b21-2e64-4ddc-991e-cae93055f798")
	documentFileID4 := uuid.MustParse("44487b21-2e64-4ddc-991e-cae93055f798")

	time1 := time.Date(2030, 10, 10, 10, 10, 10, 10, time.UTC)
	time2 := time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)
	time3 := time.Date(2010, 10, 10, 10, 10, 10, 10, time.UTC)
	time4 := time.Date(2000, 10, 10, 10, 10, 10, 10, time.UTC)

	responses := []contracts.DocumentFileResponse{
		{DocumentFileID: documentFileID1, FileSide: frontFile, CreatedAt: time1},
		{DocumentFileID: documentFileID2, FileSide: frontFile, CreatedAt: time2},
		{DocumentFileID: documentFileID3, FileSide: backFile, CreatedAt: time3},
		{DocumentFileID: documentFileID4, FileSide: frontFile, CreatedAt: time4},
	}

	expectedFrontFileResponses := []contracts.DocumentFileResponse{
		{DocumentFileID: documentFileID1, FileSide: frontFile, CreatedAt: time1},
		{DocumentFileID: documentFileID2, FileSide: frontFile, CreatedAt: time2},
		{DocumentFileID: documentFileID4, FileSide: frontFile, CreatedAt: time4},
	}

	expectedBackFileResponses := []contracts.DocumentFileResponse{
		{DocumentFileID: documentFileID3, FileSide: backFile, CreatedAt: time3},
	}

	receivedFrontFileResponses, receivedBackFileResponses := splitResponsesByFileSide(responses)

	if !reflect.DeepEqual(expectedFrontFileResponses, receivedFrontFileResponses) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expectedFrontFileResponses, receivedFrontFileResponses)
	}

	if !reflect.DeepEqual(expectedBackFileResponses, receivedBackFileResponses) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expectedBackFileResponses, receivedBackFileResponses)
	}

}

func TestSortResponsesDescending(t *testing.T) {
	documentFileID1 := uuid.MustParse("111a6ce0-9a19-43b8-9776-44ddab16d5b4")
	documentFileID2 := uuid.MustParse("222912b9-5ce1-495e-8681-be818744954a")
	documentFileID3 := uuid.MustParse("33387b21-2e64-4ddc-991e-cae93055f798")

	time1 := time.Date(2000, 10, 10, 10, 10, 10, 10, time.UTC)
	time2 := time.Date(2010, 10, 10, 10, 10, 10, 10, time.UTC)
	time3 := time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)

	responses := []contracts.DocumentFileResponse{
		{DocumentFileID: documentFileID1, CreatedAt: time1},
		{DocumentFileID: documentFileID2, CreatedAt: time2},
		{DocumentFileID: documentFileID3, CreatedAt: time3},
	}

	expected := []contracts.DocumentFileResponse{
		{DocumentFileID: documentFileID3, CreatedAt: time3},
		{DocumentFileID: documentFileID2, CreatedAt: time2},
		{DocumentFileID: documentFileID1, CreatedAt: time1},
	}

	received := sortResponsesDescending(responses)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestFindLast(t *testing.T) {
	documentFileID1 := uuid.MustParse("111a6ce0-9a19-43b8-9776-44ddab16d5b4")
	documentFileID2 := uuid.MustParse("222912b9-5ce1-495e-8681-be818744954a")
	documentFileID3 := uuid.MustParse("33387b21-2e64-4ddc-991e-cae93055f798")

	time1 := time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)
	time2 := time.Date(2010, 10, 10, 10, 10, 10, 10, time.UTC)
	time3 := time.Date(2000, 10, 10, 10, 10, 10, 10, time.UTC)

	sortedDocumentFiles := []entity.DocumentFile{
		{DocumentFileID: &documentFileID1, CreatedAt: time1},
		{DocumentFileID: &documentFileID2, CreatedAt: time2},
		{DocumentFileID: &documentFileID3, CreatedAt: time3},
	}

	expected := &sortedDocumentFiles[0]

	received := findLast(sortedDocumentFiles)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestFindLast_NoElements(t *testing.T) {
	sortedDocumentFiles := []entity.DocumentFile{}

	var expected *entity.DocumentFile = nil

	received := findLast(sortedDocumentFiles)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

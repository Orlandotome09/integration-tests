package file

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/file/http"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestGetUrl(t *testing.T) {
	client := &fileHttpClient.MockFileHttpClient{}

	service := NewFileAdapter(client)

	fileID := uuid.New()
	fileResponse := &fileHttpClient.FileResponse{FileID: fileID.String(), Url: "some url"}

	client.On("GetFileUrl", fileID.String()).Return(fileResponse, nil)

	expected := fileResponse.Url
	received, err := service.GetUrl(fileID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

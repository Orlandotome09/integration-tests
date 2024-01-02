package questionForm

import (
	questionFormClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/questionForm/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/questionForm/http/contracts"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
)

var (
	client  *questionFormClient.MockQuestionFormClient
	service interfaces.QuestionFormAdapter
)

func TestMain(m *testing.M) {
	client = &questionFormClient.MockQuestionFormClient{}
	service = New(client)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	id := uuid.New()

	resp := &contracts.QuestionFormResponse{EntityID: "174ab428-375e-48b0-9998-7484e8738304"}

	client.On("Get", id.String()).Return(resp, nil)

	expected := &entity.QuestionForm{EntityID: uuid.MustParse(resp.EntityID)}
	received, err := service.Get(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

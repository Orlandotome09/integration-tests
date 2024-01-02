package contact

import (
	contactClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/contracts"
	contactTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
)

var (
	client     *contactClient.MockContactClient
	service    interfaces.ContactAdapter
	translator contactTranslator.Translator
)

func TestMain(m *testing.M) {
	client = &contactClient.MockContactClient{}
	translator = &contactTranslator.MockContactTranslator{}
	service = NewContactAdapter(client, translator)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	id := uuid.New()

	resp := &contracts.ContactResponse{ProfileID: uuid.MustParse("174ab428-375e-48b0-9998-7484e8738304")}

	client.On("Get", id.String()).Return(resp, nil)

	expected := &entity.Contact{ProfileID: &resp.ProfileID}
	received, err := service.Get(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

package foreignAccount

import (
	foreignAccountClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http/contracts"
	foreignAccountTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"os"
	"reflect"
	"testing"
)

var (
	client     *foreignAccountClient.MockForeignAccountClient
	translator *foreignAccountTranslator.MockForeignAccountTranslator
	adapter    interfaces.ForeignAccountAdapter
)

func TestMain(m *testing.M) {
	client = &foreignAccountClient.MockForeignAccountClient{}
	translator = &foreignAccountTranslator.MockForeignAccountTranslator{}
	adapter = NewForeignAccountAdapter(client, translator)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	foreignAccountID := "1111"
	response := &contracts.ForeignAccountResponse{}
	foreignAccount := entity.ForeignAccount{}

	client.On("Get", foreignAccountID).Return(response, nil)
	translator.On("Translate", *response).Return(foreignAccount)

	expected := &foreignAccount
	received, err := adapter.Get(foreignAccountID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGet_NoContent(t *testing.T) {
	foreignAccountID := "2222"
	var response *contracts.ForeignAccountResponse = nil

	client.On("Get", foreignAccountID).Return(response, nil)

	var expected *entity.ForeignAccount
	received, err := adapter.Get(foreignAccountID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

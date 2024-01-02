package address

import (
	addressClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/contracts"
	registrationAddressesTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
)

var (
	addressclient *addressClient.MockAddressClient
	translator    *registrationAddressesTranslator.MockAddressTranslator
	service       interfaces.AddressAdapter
)

func TestMain(m *testing.M) {
	addressclient = &addressClient.MockAddressClient{}
	translator = &registrationAddressesTranslator.MockAddressTranslator{}
	service = NewAddressAdapter(addressclient, translator)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	id := "111"
	resp := &contracts.AddressResponse{ProfileID: uuid.New()}

	addressclient.On("Get", id).Return(resp, nil)

	expected := &entity.Address{ProfileID: &resp.ProfileID}

	received, err := service.Get(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestGet_NoContent(t *testing.T) {
	id := "222"
	var resp *contracts.AddressResponse = nil

	addressclient.On("Get", id).Return(resp, nil)

	var expected *entity.Address = nil

	received, err := service.Get(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestSearch(t *testing.T) {
	id := "333"
	responses := []contracts.AddressResponse{{AddressID: uuid.New()}, {AddressID: uuid.New()}}
	addresses := []entity.Address{{Street: "xxx"}, {Street: "yyy"}}

	addressclient.On("Search", id).Return(responses, nil)

	translator.On("Translate", responses).Return(addresses)

	expected := addresses

	received, err := service.Search(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

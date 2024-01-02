package legalRepresentative

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"os"
	"reflect"
	"testing"

	legalRepresentativeClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http/contracts"
	legalRepresentativeTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"github.com/google/uuid"
)

var (
	client     *legalRepresentativeClient.MockLegalRepresentativeClient
	translator *legalRepresentativeTranslator.MockLegalRepresentativeTranslator
	adapter    interfaces.LegalRepresentativeAdapter
)

func TestMain(m *testing.M) {
	client = &legalRepresentativeClient.MockLegalRepresentativeClient{}
	translator = &legalRepresentativeTranslator.MockLegalRepresentativeTranslator{}
	adapter = New(client, translator)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	id := uuid.New()

	resp := &contracts.LegalRepresentativeResponse{ProfileID: "174ab428-375e-48b0-9998-7484e8738304"}

	client.On("Get", id.String()).Return(resp, nil)

	expected := &entity2.LegalRepresentative{
		Person: entity2.Person{
			ProfileID: uuid.MustParse(resp.ProfileID),
		},
	}

	translator.On("Translate", *resp).Return(expected, nil)

	received, err := adapter.Get(id)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestSearch(t *testing.T) {

	profileID := uuid.New()

	resp := []contracts.LegalRepresentativeResponse{{ProfileID: profileID.String()}}

	client.On("Search", profileID.String()).Return(resp, nil)

	expected := []entity2.LegalRepresentative{
		{
			Person: entity2.Person{
				ProfileID: profileID,
			},
		},
	}

	translator.On("Translate", resp[0]).Return(&expected[0], nil)

	received, err := adapter.Search(profileID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

package boardOfDirectors

import (
	boardOfDirectorsClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http/contracts"
	boardOfDirectorsTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
)

var (
	client     *boardOfDirectorsClient.MockBoardOfDirectorsClient
	translator *boardOfDirectorsTranslator.MockBoardOfDirectorsTranslator
	adapter    interfaces.BoardOfDirectorsAdapter
)

func TestMain(m *testing.M) {
	client = &boardOfDirectorsClient.MockBoardOfDirectorsClient{}
	translator = &boardOfDirectorsTranslator.MockBoardOfDirectorsTranslator{}
	adapter = New(client, translator)
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {

	profileID := uuid.New()

	resp := []contracts.BoardOfDirectorsResponse{{ProfileID: profileID}}

	client.On("Search", profileID.String()).Return(resp, nil)

	expected := []entity2.Director{
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

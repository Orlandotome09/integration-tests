package foreignAccountTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := New()
	foreignAccountID := uuid.New()
	profileID := uuid.New()

	response := contracts.ForeignAccountResponse{ForeignAccountID: foreignAccountID, ProfileID: profileID}

	expected := entity.ForeignAccount{ForeignAccountID: response.ForeignAccountID, ProfileID: response.ProfileID}

	received := translator.Translate(response)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

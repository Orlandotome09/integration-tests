package addressTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate(t *testing.T) {
	addressID := uuid.New()
	profileID := uuid.New()

	addresses := []entity.Address{{
		AddressID:    &addressID,
		ProfileID:    &profileID,
		Type:         "TypeX",
		ZipCode:      "01234",
		Street:       "Rua X",
		Number:       "123",
		Complement:   "A",
		Neighborhood: "N",
		City:         "C",
		StateCode:    "A",
		CountryCode:  "BRA",
	}}

	expectedAddresses := message.Addresses{message.Address{
		Type:         addresses[0].Type,
		Street:       addresses[0].Street,
		Number:       addresses[0].Number,
		Complement:   addresses[0].Complement,
		Neighborhood: addresses[0].Neighborhood,
		City:         addresses[0].City,
		State:        addresses[0].StateCode,
		Country:      addresses[0].CountryCode,
		Code:         addresses[0].ZipCode,
	}}

	addressTranslator := addressTranslator{}

	translatedAddresses := addressTranslator.Translate(addresses)

	assert.Equal(t, expectedAddresses, translatedAddresses)
}

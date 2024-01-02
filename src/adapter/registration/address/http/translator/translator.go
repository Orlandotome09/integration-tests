package registrationAddressesTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"sort"
)

type Translator interface {
	Translate(responses []contracts.AddressResponse) []entity.Address
}

type addressesTranslator struct{}

func New() Translator {
	return &addressesTranslator{}
}

func (ref *addressesTranslator) Translate(responses []contracts.AddressResponse) []entity.Address {
	responses = sortByUpdateDateDesc(responses)

	addresses := []entity.Address{}
	for i, response := range responses {
		address := entity.Address{
			AddressID:    &responses[i].AddressID,
			ProfileID:    &responses[i].ProfileID,
			Type:         response.Type,
			ZipCode:      response.ZipCode,
			Street:       response.Street,
			Number:       response.Number,
			Complement:   response.Complement,
			Neighborhood: response.Neighborhood,
			City:         response.City,
			StateCode:    response.StateCode,
			CountryCode:  response.CountryCode,
		}
		addresses = append(addresses, address)
	}

	return addresses
}

func sortByUpdateDateDesc(responses []contracts.AddressResponse) []contracts.AddressResponse {
	sort.SliceStable(responses, func(i, j int) bool {
		return responses[i].UpdatedAt.After(responses[j].UpdatedAt)
	})

	return responses
}

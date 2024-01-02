package addressTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type AddressTranslator interface {
	Translate(addresses []entity.Address) message.Addresses
}

type addressTranslator struct{}

func NewAddressTranslator() AddressTranslator {
	return &addressTranslator{}
}

func (ref *addressTranslator) Translate(addresses []entity.Address) message.Addresses {
	messageAddresses := make([]message.Address, len(addresses))
	for i, address := range addresses {
		messageAddresses[i] = message.Address{
			Type:         address.Type,
			Street:       address.Street,
			Number:       address.Number,
			Complement:   address.Complement,
			Neighborhood: address.Neighborhood,
			City:         address.City,
			State:        address.StateCode,
			Country:      address.CountryCode,
			Code:         address.ZipCode,
		}
	}
	return messageAddresses
}

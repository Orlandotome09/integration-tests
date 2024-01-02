package address

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http"
	addressesTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type addressAdapter struct {
	addressClient       addressClient.AddressClient
	addressesTranslator addressesTranslator.Translator
}

func NewAddressAdapter(addressClient addressClient.AddressClient,
	addressesTranslator addressesTranslator.Translator,
) interfaces.AddressAdapter {
	return &addressAdapter{
		addressClient:       addressClient,
		addressesTranslator: addressesTranslator,
	}
}

func (ref *addressAdapter) Get(id string) (*entity.Address, error) {
	resp, err := ref.addressClient.Get(id)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	address := &entity.Address{ProfileID: &resp.ProfileID}

	return address, nil
}

func (ref *addressAdapter) Search(profileID string) ([]entity.Address, error) {
	responses, err := ref.addressClient.Search(profileID)
	if err != nil {
		return nil, err
	}

	addresses := ref.addressesTranslator.Translate(responses)

	return addresses, nil
}

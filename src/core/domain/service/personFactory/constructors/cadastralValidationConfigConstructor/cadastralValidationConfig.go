package cadastralValidationConfigConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"context"
	"github.com/pkg/errors"
)

type cadastralValidationConfigConstructor struct {
	cadastralValidationConfigAdapter interfaces.TemisConfigAdapter
}

func New(adapter interfaces.TemisConfigAdapter) interfaces.PersonConstructor {
	return &cadastralValidationConfigConstructor{
		cadastralValidationConfigAdapter: adapter,
	}
}

func (ref *cadastralValidationConfigConstructor) Assemble(personWrapper *entity.PersonWrapper) error {

	roleType := personWrapper.Person.RoleType
	personType := personWrapper.Person.PersonType
	partnerID := personWrapper.Person.PartnerID
	offerType := personWrapper.Person.OfferType

	cadastralValidationConfig, err := ref.cadastralValidationConfigAdapter.GetCadastralValidationConfig(context.Background(), personType, roleType, offerType, partnerID)
	if err != nil {
		return errors.WithStack(err)
	}

	if cadastralValidationConfig == nil {
		return nil
	}

	personWrapper.Person.CadastralValidationConfig = cadastralValidationConfig

	return nil
}
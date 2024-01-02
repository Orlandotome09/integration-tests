package legalRepresentative

import (
	legalRepresentativeClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http"
	legalRepresentativeTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type legalRepresentativeAdapter struct {
	client     legalRepresentativeClient.LegalRepresentativeClient
	translator legalRepresentativeTranslator.LegalRepresentativeTranslator
}

func New(client legalRepresentativeClient.LegalRepresentativeClient, translator legalRepresentativeTranslator.LegalRepresentativeTranslator) interfaces.LegalRepresentativeAdapter {
	return &legalRepresentativeAdapter{
		client:     client,
		translator: translator,
	}
}

func (ref *legalRepresentativeAdapter) Get(id uuid.UUID) (*entity.LegalRepresentative, error) {
	resp, err := ref.client.Get(id.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, nil
	}

	return ref.translator.Translate(*resp)

}

func (ref *legalRepresentativeAdapter) Search(profileID uuid.UUID) ([]entity.LegalRepresentative, error) {
	resp, err := ref.client.Search(profileID.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, nil
	}

	result := make([]entity.LegalRepresentative, len(resp))
	for index, item := range resp {
		lr, err := ref.translator.Translate(item)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result[index] = *lr
	}

	return result, nil

}

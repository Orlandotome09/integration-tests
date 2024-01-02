package boardOfDirectors

import (
	boardOfDirectorsClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http"
	boardOfDirectorsTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type boardOfDirectorsAdapter struct {
	client     boardOfDirectorsClient.BoardOfDirectorsClient
	translator boardOfDirectorsTranslator.BoardOfDirectorsTranslator
}

func New(client boardOfDirectorsClient.BoardOfDirectorsClient, translator boardOfDirectorsTranslator.BoardOfDirectorsTranslator) interfaces.BoardOfDirectorsAdapter {
	return &boardOfDirectorsAdapter{
		client:     client,
		translator: translator,
	}
}

func (ref *boardOfDirectorsAdapter) Search(profileID uuid.UUID) ([]entity.Director, error) {
	resp, err := ref.client.Search(profileID.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, nil
	}

	result := make([]entity.Director, len(resp))
	for index, item := range resp {
		director, err := ref.translator.Translate(item)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result[index] = *director
	}

	return result, nil

}

package questionForm

import (
	questionFormClient "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/questionForm/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type questionFormAdapter struct {
	client questionFormClient.QuestionFormClient
}

func New(client questionFormClient.QuestionFormClient) interfaces.QuestionFormAdapter {
	return &questionFormAdapter{
		client: client,
	}
}

func (ref *questionFormAdapter) Get(id uuid.UUID) (*entity.QuestionForm, error) {
	resp, err := ref.client.Get(id.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, nil
	}

	questionForm := &entity.QuestionForm{EntityID: uuid.MustParse(resp.EntityID)}

	return questionForm, nil
}

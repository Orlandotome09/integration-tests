package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type QuestionFormAdapter interface {
	Get(id uuid.UUID) (*entity.QuestionForm, error)
}

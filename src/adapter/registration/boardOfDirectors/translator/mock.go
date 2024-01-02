package boardOfDirectorsTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockBoardOfDirectorsTranslator struct {
	BoardOfDirectorsTranslator
	mock.Mock
}

func (ref *MockBoardOfDirectorsTranslator) Translate(response contracts.BoardOfDirectorsResponse) (*entity.Director, error) {
	args := ref.Called(response)
	return args.Get(0).(*entity.Director), args.Error(1)
}

package questionFormClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/questionForm/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockQuestionFormClient struct {
	mock.Mock
	QuestionFormClient
}

func (ref *MockQuestionFormClient) Get(id string) (*contracts.QuestionFormResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.QuestionFormResponse), args.Error(1)
}

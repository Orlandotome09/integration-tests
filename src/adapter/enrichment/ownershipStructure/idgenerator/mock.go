package idgenerator

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockIdGenerator struct {
	IdGenerator
	mock.Mock
}

func (ref *MockIdGenerator) Generate(legalEntity string, documentNumber string) uuid.UUID {
	args := ref.Called(legalEntity, documentNumber)
	return args.Get(0).(uuid.UUID)
}

package auth

import (
	interfacesAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/_interfacesAdapter"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	interfacesAdapter.AuthRepository
	mock.Mock
}

func (ref *MockAuthRepository) GetAccessToken() (string, error) {
	args := ref.Called()
	return args.Get(0).(string), args.Error(1)
}

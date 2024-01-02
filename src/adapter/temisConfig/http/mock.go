package http

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockTemisConfigHttpClient struct {
	TemisConfigHttpClient
	mock.Mock
}

func (ref *MockTemisConfigHttpClient) GetCadastralValidationConfig(ctx context.Context, personType, roleType, offerType string) (CadastralValidationConfigResponse, error) {
	args := ref.Called(ctx, personType, roleType, offerType)
	return args.Get(0).(CadastralValidationConfigResponse), args.Error(1)
}

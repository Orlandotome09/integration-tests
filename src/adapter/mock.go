package adapter

import (
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	HttpClient
	mock.Mock
}

func (ref *MockHttpClient) Host() string {
	args := ref.Called()
	return args.Get(0).(string)
}

func (ref *MockHttpClient) Get(path string, params string, headers map[string]string) ([]byte, error) {
	args := ref.Called(path, params, headers)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (ref *MockHttpClient) Search(path string, params string) ([]byte, error) {
	args := ref.Called(path, params)
	return args.Get(0).([]byte), args.Error(1)
}

func (ref *MockHttpClient) Post(path string, body interface{}) ([]byte, error) {
	args := ref.Called(path, body)
	return args.Get(0).([]byte), args.Error(1)
}

func (ref *MockHttpClient) SetResponseInterface(ifc interface{}) {
	ref.Called(ifc)
}

func (ref *MockHttpClient) SetErrorInterface(err ErrorInterface) {
	ref.Called(err)
}

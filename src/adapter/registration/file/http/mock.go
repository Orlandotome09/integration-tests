package fileHttpClient

import "github.com/stretchr/testify/mock"

type MockFileHttpClient struct {
	FileHttpClient
	mock.Mock
}

func (ref *MockFileHttpClient) GetFileUrl(fileID string) (*FileResponse, error) {
	args := ref.Called(fileID)
	return args.Get(0).(*FileResponse), args.Error(1)
}

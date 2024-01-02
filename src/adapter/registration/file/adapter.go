package file

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/file/http"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type fileServiceUrl struct {
	fileHttpClient fileHttpClient.FileHttpClient
}

func NewFileAdapter(fileHttpClient fileHttpClient.FileHttpClient) interfaces.FileAdapter {
	return &fileServiceUrl{
		fileHttpClient: fileHttpClient,
	}
}

func (ref *fileServiceUrl) GetUrl(fileID uuid.UUID) (string, error) {
	resp, err := ref.fileHttpClient.GetFileUrl(fileID.String())
	if err != nil {
		return "", errors.WithStack(err)
	}

	return resp.Url, nil
}

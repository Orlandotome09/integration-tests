package fileHttpClient

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	filePath = "/file/"
)

type FileHttpClient interface {
	GetFileUrl(fileID string) (*FileResponse, error)
}

type fileHttpClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) FileHttpClient {
	return &fileHttpClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *fileHttpClient) GetFileUrl(fileID string) (*FileResponse, error) {
	uri := ref.host + filePath + fileID + "/url"

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting file url. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := &FileResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

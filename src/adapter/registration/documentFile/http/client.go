package documentFileClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	getPath    = "/document_file/"
	searchPath = "/document/"
)

type DocumentFileClient interface {
	Get(documentFileID string) (*contracts.DocumentFileResponse, error)
	SearchByDocumentId(id string) ([]contracts.DocumentFileResponse, error)
}

type documentFileClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) DocumentFileClient {
	return &documentFileClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *documentFileClient) Get(documentFileID string) (*contracts.DocumentFileResponse, error) {
	uri := ref.host + getPath + documentFileID

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {

		str := fmt.Sprintf("Error requesting document file. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := &contracts.DocumentFileResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func (ref *documentFileClient) SearchByDocumentId(id string) ([]contracts.DocumentFileResponse, error) {
	uri := ref.host + searchPath + id + "/files"

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting document file. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := []contracts.DocumentFileResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

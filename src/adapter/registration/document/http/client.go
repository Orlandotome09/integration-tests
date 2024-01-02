package documentClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/document/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	getPath    = "/document/"
	searchPath = "/documents"
)

type DocumentClient interface {
	Get(id string) (*contracts.DocumentResponse, error)
	SearchByEntityID(id string) ([]contracts.DocumentResponse, error)
}

type documentClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) DocumentClient {
	return &documentClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *documentClient) Get(id string) (*contracts.DocumentResponse, error) {
	uri := ref.host + getPath + id

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {

		str := fmt.Sprintf("Error requesting address. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := &contracts.DocumentResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func (ref *documentClient) SearchByEntityID(id string) ([]contracts.DocumentResponse, error) {
	uri := ref.host + searchPath

	req, _ := http.NewRequest("GET", uri, nil)
	params := req.URL.Query()
	params.Add("entity_id", id)
	req.URL.RawQuery = params.Encode()

	resp, err := ref.webClient.Do(req)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK {
		str := fmt.Sprintf("Error requesting documents. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		logrus.Errorf("[DocumentClient] %v", str)
		return nil, errors.WithStack(errors.New(str))
	}

	response := []contracts.DocumentResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}
	return response, nil
}

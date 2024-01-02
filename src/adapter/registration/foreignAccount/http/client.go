package foreignAccountClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/foreignAccount/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	getPath = "/foreign-accounts/"
)

type ForeignAccountClient interface {
	Get(id string) (*contracts.ForeignAccountResponse, error)
}

type foreignAccountClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) ForeignAccountClient {
	return &foreignAccountClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *foreignAccountClient) Get(id string) (*contracts.ForeignAccountResponse, error) {
	uri := ref.host + getPath + id
	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting foreign account. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := contracts.ForeignAccountResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

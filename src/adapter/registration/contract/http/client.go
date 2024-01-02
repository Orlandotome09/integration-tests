package contractClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http/contracts"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const (
	getContractPath = "/contract/"
)

type ContractClient interface {
	Get(id *uuid.UUID) (*contracts.GetContractResponse, bool, error)
}

type contractHttpClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) ContractClient {
	return &contractHttpClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *contractHttpClient) Get(id *uuid.UUID) (*contracts.GetContractResponse, bool, error) {
	uri := ref.host + getContractPath + id.String()

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, false, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, false, nil
		}
		str := fmt.Sprintf("Error requesting contract. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, false, errors.New(str)
	}

	response := &contracts.GetContractResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, false, err
	}

	return response, true, nil
}

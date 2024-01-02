package addressClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/address/http/contracts"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	getPath    = "/address/"
	searchPath = "/addresses"
)

type AddressClient interface {
	Get(id string) (*contracts.AddressResponse, error)
	Search(profileID string) ([]contracts.AddressResponse, error)
}

type addressClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) AddressClient {
	return &addressClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *addressClient) Get(id string) (*contracts.AddressResponse, error) {
	uri := ref.host + getPath + id

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {

		str := fmt.Sprintf("Error requesting address. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.New(str)
	}

	response := &contracts.AddressResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (ref *addressClient) Search(profileID string) ([]contracts.AddressResponse, error) {
	uri := ref.host + searchPath

	req, _ := http.NewRequest("GET", uri, nil)
	params := req.URL.Query()
	params.Add("profile_id", profileID)
	req.URL.RawQuery = params.Encode()

	resp, err := ref.webClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting address. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.New(str)
	}

	response := []contracts.AddressResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

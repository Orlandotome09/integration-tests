package accountClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account/http/contracts"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	getAccountPath            = "/accounts/"
	createInternalAccountPath = "/accounts/internal"
)

type AccountClient interface {
	Get(id string) (*contracts.AccountResponse, error)
	Find(profileID string) ([]contracts.AccountResponse, error)
	CreateInternal(request *contracts.CreateAccountRequest) (*contracts.AccountResponse, error)
}

type accountHttpClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) AccountClient {
	return &accountHttpClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *accountHttpClient) Get(id string) (*contracts.AccountResponse, error) {
	uri := ref.host + getAccountPath + id

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 204 {
			return nil, nil
		}
		str := fmt.Sprintf("Error requesting account. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := &contracts.AccountResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (ref *accountHttpClient) Find(profileID string) ([]contracts.AccountResponse, error) {
	uri := ref.host + "/profile/" + profileID + "/accounts"

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting account. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	var response []contracts.AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (ref *accountHttpClient) CreateInternal(request *contracts.CreateAccountRequest) (*contracts.AccountResponse, error) {
	uri := ref.host + createInternalAccountPath
	body, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(body))

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 {
		errorsResponse := contracts.ErrorsResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&errorsResponse); err != nil {
			return nil, errors.WithStack(err)

		}
		return nil, errors.WithStack(errors.New(strings.Join(errorsResponse.Error, ",")))
	}

	if resp.StatusCode != 201 && resp.StatusCode != 409 {
		errorResponse := contracts.ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.WithStack(errors.New(errorResponse.Error))
	}

	response := &contracts.AccountResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

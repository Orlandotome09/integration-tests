package legalRepresentativeClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	getPath    = "/legal-representative/"
	searchPath = "/legal-representatives"
)

type LegalRepresentativeClient interface {
	Get(id string) (*contracts.LegalRepresentativeResponse, error)
	Search(profileID string) ([]contracts.LegalRepresentativeResponse, error)
}

type legalRepresentativeClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) LegalRepresentativeClient {
	return &legalRepresentativeClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *legalRepresentativeClient) Get(id string) (*contracts.LegalRepresentativeResponse, error) {
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
		str := fmt.Sprintf("Error requesting legal representative. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := contracts.LegalRepresentativeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (ref *legalRepresentativeClient) Search(profileID string) ([]contracts.LegalRepresentativeResponse, error) {
	uri := ref.host + searchPath + "?profile_id=" + profileID

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting legal representative. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := make([]contracts.LegalRepresentativeResponse, 0)
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

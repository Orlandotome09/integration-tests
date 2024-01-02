package contactClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contact/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	getPath    = "/contact/"
	searchPath = "/contacts"
)

type ContactClient interface {
	Get(id string) (*contracts.ContactResponse, error)
	Search(profileID string) ([]contracts.ContactResponse, error)
}

type contactClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) ContactClient {
	return &contactClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *contactClient) Get(id string) (*contracts.ContactResponse, error) {
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
		str := fmt.Sprintf("Error requesting contact. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := contracts.ContactResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

func (ref *contactClient) Search(profileID string) ([]contracts.ContactResponse, error) {
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
		str := fmt.Sprintf("Error requesting contacts. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.New(str)
	}

	response := []contracts.ContactResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

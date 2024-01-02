package boardOfDirectorsClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	searchPath = "/directors"
)

type BoardOfDirectorsClient interface {
	Search(profileID string) ([]contracts.BoardOfDirectorsResponse, error)
}

type boardOfDirectorsClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) BoardOfDirectorsClient {
	return &boardOfDirectorsClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *boardOfDirectorsClient) Search(profileID string) ([]contracts.BoardOfDirectorsResponse, error) {
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

	response := make([]contracts.BoardOfDirectorsResponse, 0)
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

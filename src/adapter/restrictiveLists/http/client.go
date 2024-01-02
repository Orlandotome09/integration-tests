package restrictiveListsHttpClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	searchInternaListPath = "/internal-list"
	searchPepList         = "/pep"
)

type RestrictiveListsHttpClient interface {
	SearchInternalList(documentFilter string, nameFilter string) (InternalListResponse, error)
	SearchPepList(documentNumber string) (*PepResponse, error)
}

type restrictiveListsHttpClient struct {
	webClient *http.Client
	host      string
}

func NewRestrictiveListHttpClient(webClient *http.Client, host string) RestrictiveListsHttpClient {
	return &restrictiveListsHttpClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *restrictiveListsHttpClient) SearchInternalList(documentFilter string, nameFilter string) (InternalListResponse, error) {
	uri := ref.host + searchInternaListPath

	req, _ := http.NewRequest("GET", uri, nil)
	params := req.URL.Query()
	params.Add("document_number", documentFilter)
	params.Add("full_name", nameFilter)
	req.URL.RawQuery = params.Encode()

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting internal lists. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.New(str)
	}

	response := InternalListResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (ref *restrictiveListsHttpClient) SearchPepList(documentNumber string) (*PepResponse, error) {
	uri := ref.host + searchPepList + "/" + documentNumber

	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		str := fmt.Sprintf("Error requesting pep lists. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.New(str)
	}

	response := &PepResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

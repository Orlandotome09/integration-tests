package ownershipStructureHttpClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	GetOwnershipStructurePath = "/ownership-structure/"
)

type OwnershipStructureClient interface {
	GetOwnershipStructureForProfile(profileID string) (*contracts.OwnershipStructureResponse, error)
}

type ownershipStructureClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) OwnershipStructureClient {
	return &ownershipStructureClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *ownershipStructureClient) GetOwnershipStructureForProfile(profileID string) (*contracts.OwnershipStructureResponse, error) {
	uri := ref.host + GetOwnershipStructurePath + profileID
	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 204 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting ownership structure. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := contracts.OwnershipStructureResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

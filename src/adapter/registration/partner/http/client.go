package partnerClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	getPath = "/partner/"
)

type PartnerClient interface {
	Get(id string) (*contracts.PartnerResponse, error)
}

type partnerClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) PartnerClient {
	return &partnerClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *partnerClient) Get(id string) (*contracts.PartnerResponse, error) {
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
		str := fmt.Sprintf("Error requesting partner. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := contracts.PartnerResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

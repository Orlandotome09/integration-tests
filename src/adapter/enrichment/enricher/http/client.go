package enricherClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/enricher/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	GetEnricherPath = "/enrich/"
)

type EnricherClient interface {
	GetEnrichedPerson(request enricherContracts.EnricherRequest, documentNumber string) (*enricherContracts.EnricherResponse, error)
}

type enricherClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) EnricherClient {
	return &enricherClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *enricherClient) GetEnrichedPerson(request enricherContracts.EnricherRequest, documentNumber string) (*enricherContracts.EnricherResponse, error) {

	documentNumberNormalized := core.NormalizeDocument(documentNumber)
	uri := ref.host + GetEnricherPath + documentNumberNormalized
	req, _ := http.NewRequest("GET", uri, nil)

	params := req.URL.Query()
	params.Add("profile_id", request.ProfileID)
	params.Add("person_type", request.PersonType)
	params.Add("offer_type", request.OfferType)
	params.Add("partner_id", request.PartnerID)
	params.Add("role_type", request.RoleType)
	req.URL.RawQuery = params.Encode()

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	logrus.Infof("[GetEnrichedPerson] Response from Enrichment %v and url %s", resp.StatusCode, uri)

	if resp.StatusCode == 204 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting enriched person. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := enricherContracts.EnricherResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

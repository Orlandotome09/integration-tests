package ownershipStructureClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/ownershipStructure/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	GetOwnershipStructurePath = "/ownership-structure/"
)

type OwnershipStructureClient interface {
	GetOwnershipStructure(legalEntityID, offerType, partnerID string) (*contracts.OwnershipStructureResponse, error)
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

func (ref *ownershipStructureClient) GetOwnershipStructure(legalEntityID, offerType, partnerID string) (*contracts.OwnershipStructureResponse, error) {

	documentNumber := core.NormalizeDocument(legalEntityID)
	uri := ref.host + GetOwnershipStructurePath + documentNumber
	req, _ := http.NewRequest("GET", uri, nil)

	req.Header.Add("Offer-Type", offerType)
	req.Header.Add("Partner-Id", partnerID)

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

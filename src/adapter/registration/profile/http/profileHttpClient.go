package profileHttpClient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/profile/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

const (
	getProfilePath                  = "/profile/"
	createProfilePath               = "/profile"
	updateProfilePath               = "/profile/"
	findByDocumentNumberProfilePath = "/profiles/find_by_document_number"
	contentType                     = "application/json; charset=utf-8"
)

type ProfileHttpClient interface {
	Get(profileID uuid.UUID) (*contracts.ProfileResponse, error)
	Create(request contracts.ProfileRequest) (*contracts.ProfileResponse, error)
	FindByDocumentNumber(roleType values.RoleType, documentNumber string, partnerID string, parentID *uuid.UUID) (*contracts.ProfileResponse, error)
}

type profileHttpClient struct {
	webClient *http.Client
	host      string
}

func NewProfileHttpClient(webClient *http.Client, host string) ProfileHttpClient {
	return &profileHttpClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *profileHttpClient) Get(profileID uuid.UUID) (*contracts.ProfileResponse, error) {
	uri := ref.host + getProfilePath + profileID.String()

	resp, err := ref.webClient.Get(uri)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, errors.Errorf("Error requesting profile get. Status Code: %v. Uri: %v. Body: %s", resp.StatusCode, uri, bodyString)
	}

	response := &contracts.ProfileResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func (ref *profileHttpClient) Create(request contracts.ProfileRequest) (*contracts.ProfileResponse, error) {
	uri := ref.host + createProfilePath

	body, _ := json.Marshal(request)

	resp, err := ref.webClient.Post(uri, contentType, bytes.NewBuffer(body))

	if err != nil {

		return nil, err

	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, errors.Errorf("Error requesting profile creation. Status Code: %v. Uri: %v. Body: %s", resp.StatusCode, uri, bodyString)
	}

	response := &contracts.ProfileResponse{}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func (ref *profileHttpClient) FindByDocumentNumber(roleType values.RoleType, documentNumber string, partnerID string, parentID *uuid.UUID) (*contracts.ProfileResponse, error) {
	uri := ref.host + findByDocumentNumberProfilePath

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	params := req.URL.Query()
	params.Add("role_type", string(roleType))
	params.Add("document_number", documentNumber)
	params.Add("partner_id", partnerID)
	if parentID != nil {
		params.Add("parent_id", parentID.String())
	}

	req.URL.RawQuery = params.Encode()

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, errors.Errorf("Error requesting profile find_by_document_number. Status Code: %v. Uri: %v. Body: %s", resp.StatusCode, uri, bodyString)
	}

	response := &contracts.ProfileResponse{}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

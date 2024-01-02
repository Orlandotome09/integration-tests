package http

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	cadastralValidationConfigPath = "/cadastral-validation-configs"
	cacheTTL                      = time.Minute * 120
)

type TemisConfigHttpClient interface {
	GetCadastralValidationConfig(ctx context.Context, personType, roleType, offerType string) (CadastralValidationConfigResponse, error)
}

type temisConfigHttpClient struct {
	httpClient *http.Client
	host       string
	cache      _interfaces.Cache
}

func NewTemisConfigHttpClient(
	httpClient *http.Client,
	host string,
	cache _interfaces.Cache,
) TemisConfigHttpClient {
	return &temisConfigHttpClient{
		httpClient: httpClient,
		host:       host,
		cache:      cache,
	}
}

func (ref *temisConfigHttpClient) GetCadastralValidationConfig(ctx context.Context,
	personType, roleType, offerType string) (CadastralValidationConfigResponse, error) {
	key := personType + roleType + offerType
	item := ref.cache.Get(key)
	if item != nil {
		configResponse, ok := item.(CadastralValidationConfigResponse)
		if ok {
			logrus.Infof("[TemisConfigHttpClient] found cadastral validation config in cache: %v", &configResponse)
			return configResponse, nil
		}
	}

	response, err := ref.getCadastralValidationConfig(ctx, personType, roleType, offerType)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 {
		return nil, nil
	}
	ref.cache.Save(key, response, cacheTTL)

	return response, nil
}

func (ref *temisConfigHttpClient) getCadastralValidationConfig(ctx context.Context,
	personType, roleType, offerType string) (CadastralValidationConfigResponse, error) {
	uri := ref.host + cadastralValidationConfigPath

	req, _ := http.NewRequest("GET", uri, nil)
	params := req.URL.Query()

	params.Add("person_type", personType)
	params.Add("role_type", roleType)
	params.Add("offer", offerType)

	req.URL.RawQuery = params.Encode()

	resp, err := ref.httpClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != 200 {
		var errBody []byte
		resp.Body.Read(errBody)
		msg := fmt.Sprintf("Error requesting cadastral validation config. Status Code: %v. Uri: %v. Response: %v", resp.StatusCode, uri, string(errBody))
		return nil, errors.WithStack(errors.New(msg))
	}

	var response CadastralValidationConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

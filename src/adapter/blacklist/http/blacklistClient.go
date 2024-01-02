package blacklistClient

import (
	interfacesAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/_interfacesAdapter"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/auth"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http/dto"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const (
	BlacklistServiceURI = "/blacklist"
)

type BlackListHttpClient interface {
	Search(documentNumber, partnerId string) (dto.BlacklistResponse, bool, error)
}

type blacklistHttpClient struct {
	webClient      *http.Client
	host           string
	authRepository interfacesAdapter.AuthRepository
}

func New(webClient *http.Client, host string, authRepository interfacesAdapter.AuthRepository) BlackListHttpClient {
	return &blacklistHttpClient{
		webClient:      webClient,
		host:           host,
		authRepository: authRepository,
	}
}

func (ref *blacklistHttpClient) Search(documentNumber, partnerId string) (dto.BlacklistResponse, bool, error) {

	params := documentNumber + "?partnerId=" + partnerId
	url := auth.JoinURL(ref.host, BlacklistServiceURI)
	url = auth.JoinURL(url, params)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto.BlacklistResponse{}, false, errors.WithStack(err)
	}

	if ref.authRepository != nil {
		authToken, err := ref.authRepository.GetAccessToken()
		if err != nil {
			logrus.WithError(err).Errorf("[BlackListAdapter]There was an error getting access token: %s", err)
			return dto.BlacklistResponse{}, false, errors.WithStack(err)
		}

		request.Header.Add("Authorization", "Bearer "+authToken)
	}

	resp, err := ref.webClient.Do(request)
	if err != nil {
		logrus.WithError(err).Errorf("[BlackListAdapter]There was an error in blacklist integration: %s. URL: %v", err, url)
		return dto.BlacklistResponse{}, false, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		invalidStatusCodeErr := fmt.Errorf("[BlackListAdapter]Blacklist status code is invalid: %v. Body response: %v. URL: %v", resp.StatusCode, bodyString, url)
		logrus.Error(invalidStatusCodeErr)
		return dto.BlacklistResponse{}, false, invalidStatusCodeErr
	}

	var response = dto.BlacklistResponse{}

	if err := json.Unmarshal([]byte(bodyString), &response); err != nil {
		logrus.Errorf("There was an error decoding blacklist response: %+v", errors.WithStack(err))
		return dto.BlacklistResponse{}, false, err
	}

	return response, true, nil
}

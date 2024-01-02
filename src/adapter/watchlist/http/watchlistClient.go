package watchlistClient

import (
	interfacesAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/_interfacesAdapter"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/auth"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http/dto"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	WatchlistServiceURI = "/watchlist"
)

type WatchListClient interface {
	SearchIndividual(documentNumber, firstName, lastName, fullName, countryCode string, birthYear int) ([]dto.WatchlistResponse, error)
	SearchCompany(legalName, contryCode string) ([]dto.WatchlistResponse, error)
}

type watchlistHttpClient struct {
	webClient      *http.Client
	host           string
	authRepository interfacesAdapter.AuthRepository
}

func New(webClient *http.Client, host string, authRepository interfacesAdapter.AuthRepository) WatchListClient {
	return &watchlistHttpClient{
		webClient:      webClient,
		host:           host,
		authRepository: authRepository,
	}
}

func (ref *watchlistHttpClient) SearchIndividual(documentNumber, firstName, lastName, fullName, countryCode string, birthYear int) ([]dto.WatchlistResponse, error) {

	params := "firstName=" + firstName + "&lastName=" + lastName + "&birthYear=" + strconv.Itoa(birthYear) + "&fullName=" + fullName + "&countryCode=" + countryCode
	watchlistUrl := auth.JoinURL(ref.host, WatchlistServiceURI) + "/" + documentNumber + "?" + url.PathEscape(params)

	return ref.search(watchlistUrl)
}

func (ref *watchlistHttpClient) SearchCompany(legalName, countryCode string) ([]dto.WatchlistResponse, error) {

	params := "companyName=" + strings.ReplaceAll(legalName, "/", "") + "&countryCode=" + countryCode
	watchlistUrl := auth.JoinURL(ref.host, WatchlistServiceURI) + "/company" + "?" + url.PathEscape(params)

	return ref.search(watchlistUrl)
}

func (ref *watchlistHttpClient) search(url string) ([]dto.WatchlistResponse, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if ref.authRepository != nil {
		req, err = ref.addAuth(*req)
		if err != nil {
			return nil, err
		}
	}

	resp, err := ref.webClient.Do(req)
	if err != nil {
		logrus.WithField("URL", url).WithError(err).Errorf("There was an error in watchlist integration: %s", err)
		return nil, errors.WithStack(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	if resp.StatusCode != http.StatusOK {
		invalidStatusCodeErr := errors.Errorf("Watchlist status code is invalid: %d - Body: %s", resp.StatusCode, bodyString)
		logrus.WithField("URL", url).Error(invalidStatusCodeErr)
		return nil, invalidStatusCodeErr
	}

	defer resp.Body.Close()

	var response = make([]dto.WatchlistResponse, 0)

	if err := json.Unmarshal([]byte(bodyString), &response); err != nil {
		logrus.WithField("URL", url).Errorf("There was an error decoding watchlist response: %+v", errors.WithStack(err))
		return nil, err
	}

	return response, nil
}

func (ref *watchlistHttpClient) addAuth(request http.Request) (*http.Request, error) {
	authToken, err := ref.authRepository.GetAccessToken()
	if err != nil {
		logrus.WithError(err).Errorf("There was an error to get access token: %s", err)
		return nil, errors.WithStack(err)
	}

	request.Header.Add("Authorization", "Bearer "+authToken)

	return &request, nil
}

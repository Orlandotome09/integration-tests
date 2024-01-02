package auth

import (
	interfacesAdapter "bitbucket.org/bexstech/temis-compliance/src/adapter/_interfacesAdapter"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type authRepository struct {
	interfacesAdapter.AuthRepository
	authCache authCache
	Host      string
	Oauth     OAuth
}

func NewAuthRepository(host string, oAuth OAuth) interfacesAdapter.AuthRepository {
	return &authRepository{
		authCache: authCache{},
		Host:      host,
		Oauth:     oAuth,
	}
}

type OAuth struct {
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Audience     string `json:"audience,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
}

type tokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type authCache struct {
	sy    sync.Mutex
	cache *tokenResult
}

func (ref *authRepository) GetAccessToken() (string, error) {
	nowInSeconds := time.Now().Unix()
	if ref.authCache.cache != nil && ref.authCache.cache.ExpiresIn > nowInSeconds {
		return ref.authCache.cache.AccessToken, nil
	}

	tokenResult, err := ref.requestAccessToken(ref.Oauth)
	if err != nil {
		return "", err
	}

	ref.authCache.addToCache(*tokenResult)

	return tokenResult.AccessToken, err
}

func (ref *authRepository) requestAccessToken(auth OAuth) (*tokenResult, error) {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(auth); err != nil {
		return nil, fmt.Errorf("json.Encode: %+v", err)
	}

	url := JoinURL(ref.Host, "/oauth/token")
	resp, err := http.Post(url, "application/json", buf)

	if err != nil {
		logrus.WithError(err).Errorf("There was an error in auth0 integration: %s", err)
		return nil, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()

	tokenResult := &tokenResult{}

	if err := json.Unmarshal([]byte(bodyString), tokenResult); err != nil {
		logrus.Errorf("There was an error decoding auth0 response: %+v", errors.WithStack(err))
		return nil, err
	}

	return tokenResult, nil
}

func (ref *authCache) cleanFunction(key string) {
	ref.sy.Lock()
	ref.cache = nil
	ref.sy.Unlock()
}

func (ref *authCache) addToCache(token tokenResult) {
	ref.sy.Lock()
	defer ref.sy.Unlock()

	newExpiresIn := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
	token.ExpiresIn = newExpiresIn.Unix()
	ref.cache = &token
}

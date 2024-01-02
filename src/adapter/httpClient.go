package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type HttpClient interface {
	Get(path string, params string, headers map[string]string) ([]byte, error)
	Search(path string, params string) ([]byte, error)
	Post(path string, body interface{}) ([]byte, error)
}

type httpClient struct {
	webClient *http.Client
	host      string
}

func NewHttpClient(webClient *http.Client, host string) HttpClient {
	return &httpClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *httpClient) Get(path string, params string, headers map[string]string) ([]byte, error) {
	uri := ref.host + path + params

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 || resp.StatusCode == 204 {
		return nil, nil
	}
	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != 200 {
		if len(bodyResponse) > 0 {
			return nil, errors.New(fmt.Sprintf("Error making get http request. Status Code: %v. Uri: %v. Response Body: %v",
				resp.StatusCode, uri, string(bodyResponse)))
		}
		return nil, errors.New(fmt.Sprintf("Error making get http request. Status Code: %v. Uri: %v", resp.StatusCode, uri))
	}

	return bodyResponse, nil

}

func (ref *httpClient) Search(path string, params string) (bodyResponse []byte, err error) {
	return nil, nil
}

func (ref *httpClient) Post(path string, body interface{}) (bodyResponse []byte, err error) {
	uri := ref.host + path

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 || resp.StatusCode == 204 {
		return nil, nil
	}
	bodyResponse, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch resp.StatusCode {
	case 200, 201, 409:
		return
	}

	if len(bodyResponse) > 0 {
		return nil, errors.New(fmt.Sprintf("Error making get http request. Status Code: %v. Uri: %v. Response Body: %v",
			resp.StatusCode, uri, string(bodyResponse)))
	}
	return nil, errors.New(fmt.Sprintf("Error making post http request. Status Code: %v. Uri: %v", resp.StatusCode, uri))

}

package questionFormClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/questionForm/http/contracts"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

var (
	GetPath = "/question-form/"
)

type QuestionFormClient interface {
	Get(id string) (*contracts.QuestionFormResponse, error)
}

type questionFormClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) QuestionFormClient {
	return &questionFormClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *questionFormClient) Get(id string) (*contracts.QuestionFormResponse, error) {
	uri := ref.host + GetPath + id
	req, _ := http.NewRequest("GET", uri, nil)

	resp, err := ref.webClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting question form. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.WithStack(errors.New(str))
	}

	response := contracts.QuestionFormResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.WithStack(err)
	}

	return &response, nil
}

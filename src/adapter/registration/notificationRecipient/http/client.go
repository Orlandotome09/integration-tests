package notificationRecipientClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/contracts"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	searchPath = "/notification-recipients"
)

type NotificationRecipientClient interface {
	Search(profileID string) ([]contracts.NotificationRecipientResponse, error)
}

type notificationRecipientClient struct {
	webClient *http.Client
	host      string
}

func New(webClient *http.Client, host string) NotificationRecipientClient {
	return &notificationRecipientClient{
		webClient: webClient,
		host:      host,
	}
}

func (ref *notificationRecipientClient) Search(profileID string) ([]contracts.NotificationRecipientResponse, error) {
	uri := ref.host + searchPath

	req, _ := http.NewRequest("GET", uri, nil)
	params := req.URL.Query()
	params.Add("profile_id", profileID)
	req.URL.RawQuery = params.Encode()

	resp, err := ref.webClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		str := fmt.Sprintf("Error requesting notificationRecipient. Status Code: %v. Uri: %v", resp.StatusCode, uri)
		return nil, errors.New(str)
	}

	response := []contracts.NotificationRecipientResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

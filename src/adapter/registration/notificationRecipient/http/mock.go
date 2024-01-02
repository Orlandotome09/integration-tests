package notificationRecipientClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/notificationRecipient/http/contracts"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRecipientClient struct {
	NotificationRecipientClient
	mock.Mock
}

func (ref *MockNotificationRecipientClient) Get(id string) (*contracts.NotificationRecipientResponse, error) {
	args := ref.Called(id)
	return args.Get(0).(*contracts.NotificationRecipientResponse), args.Error(1)
}

func (ref *MockNotificationRecipientClient) Search(profileID string) ([]contracts.NotificationRecipientResponse, error) {
	args := ref.Called(profileID)
	return args.Get(0).([]contracts.NotificationRecipientResponse), args.Error(1)
}

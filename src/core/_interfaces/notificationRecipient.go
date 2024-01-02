package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type NotificationRecipientAdapter interface {
	Search(profileID string) ([]entity.NotificationRecipient, error)
}

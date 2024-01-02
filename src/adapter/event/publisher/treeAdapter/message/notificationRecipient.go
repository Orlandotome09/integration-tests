package message

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type NotificationRecipient struct {
	NotificationType values.NotificationType `json:"notification_type"`
	EmailTo          string                  `json:"email_to"`
	CopyEmail        string                  `json:"copy_email"`
}

type NotificationRecipients []NotificationRecipient

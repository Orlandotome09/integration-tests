package contract

type NotificationRecipient struct {
	NotificationType string `json:"notification_type"`
	EmailTo          string `json:"email_to"`
	CopyEmail        string `json:"copy_email"`
}

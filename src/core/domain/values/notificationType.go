package values

import (
	"fmt"
	"github.com/pkg/errors"
)

type NotificationType string

const (
	NotificationTypePostWarnings           NotificationType = "POST_WARNINGS"
	NotificationTypeExchangeContract       NotificationType = "EXCHANGE_CONTRACT"
	NotificationTypeSendPendingDocuments   NotificationType = "SEND_PENDING_DOCUMENTS"
	NotificationTypeSentOP                 NotificationType = "SENT_OP"
	NotificationTypeReceivedOP             NotificationType = "RECEIVED_OP"
	NotificationTypeCancelExchangeRegister NotificationType = "CANCEL_EXCHANGE_REGISTER"
)

var ValidNotificationTypes = map[string]NotificationType{
	NotificationTypePostWarnings.ToString():           NotificationTypePostWarnings,
	NotificationTypeExchangeContract.ToString():       NotificationTypeExchangeContract,
	NotificationTypeSendPendingDocuments.ToString():   NotificationTypeSendPendingDocuments,
	NotificationTypeSentOP.ToString():                 NotificationTypeSentOP,
	NotificationTypeReceivedOP.ToString():             NotificationTypeReceivedOP,
	NotificationTypeCancelExchangeRegister.ToString(): NotificationTypeCancelExchangeRegister,
}

func (notificationType *NotificationType) Validate() error {
	value := notificationType.ToString()
	if _, in := ValidNotificationTypes[value]; !in {
		return errors.New(fmt.Sprintf("%s is an invalid notification type", value))
	}
	return nil
}

func (notificationType NotificationType) ToString() string {
	return string(notificationType)
}

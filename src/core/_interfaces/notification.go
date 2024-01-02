package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type NotificationService interface {
	SendNotification(event *values.Event) (*entity.State, error)
}

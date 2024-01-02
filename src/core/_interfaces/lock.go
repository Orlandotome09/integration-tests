package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type LockRepository interface {
	BeginTransaction(event values.Event, processor func(event values.Event, timeout *bool) error) error
}

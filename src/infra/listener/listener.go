package listener

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type Processor func(event *values.Event) error

type Listener interface {
	Listen(processor Processor) error
}

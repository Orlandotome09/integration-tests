package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type EventTranslator interface {
	Translate(message []byte) (*values.Event, error)
}

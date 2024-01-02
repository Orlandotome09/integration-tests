package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type LimitPublisher interface {
	Send(profile entity2.Profile, state entity2.State) error
}

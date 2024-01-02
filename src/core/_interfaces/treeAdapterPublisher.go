package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type TreeAdapterPublisher interface {
	Send(profile entity2.Profile, accounts []entity2.Account, addresses []entity2.Address) error
}

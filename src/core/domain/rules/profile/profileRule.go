package profile

import (
	entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ProfileRule struct {
	entity.Rule
	profile entity.Profile
}

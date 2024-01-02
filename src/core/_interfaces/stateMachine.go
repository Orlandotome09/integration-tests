package _interfaces

import "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"

type StateMachine interface {
	CalculateState(state entity.State) *entity.State
}

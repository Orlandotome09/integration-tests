package _interfaces

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type StateRepository interface {
	Get(entityID uuid.UUID) (*entity.State, bool, error)
	Create(entityID uuid.UUID, engineName values.EngineName) (*entity.State, error)
	Save(state entity.State, requestDate time.Time, executionTime time.Time) (*entity.State, error)
	Search(request entity.SearchProfileStateRequest) ([]entity.State, int64, error)
	SearchContractStates(request entity.SearchProfileStateRequest) ([]entity.State, int64, error)
	FindByProfileID(profileID uuid.UUID, engine values.EngineName, result values.Result) ([]entity.State, error)
}

type StateService interface {
	Get(entityID uuid.UUID) (*entity.State, bool, error)
	Create(entityID uuid.UUID, engineName values.EngineName) (*entity.State, error)
	Update(state entity.State, requestDate time.Time, executionTime time.Time) (err error)
	SearchProfileStates(request entity.SearchProfileStateRequest) (*entity.ProfileStateList, error)
	SearchContractStates(request entity.SearchProfileStateRequest) (*entity.ProfileStateList, error)
	FindByProfileID(profileID uuid.UUID, engine values.EngineName, result values.Result) ([]entity.State, error)
	Resync(integrationIds ...string) ([]string, error)
	Reprocess(egineName string, integrationIds ...string) ([]string, error)
}

type StateManager interface {
	GetOrCreateState(eventDate time.Time, entityId uuid.UUID,
		engineName values.EngineName) (state *entity.State, shouldIgnore bool, err error)
	UpdateState(state *entity.State, requestDate time.Time, executionTime time.Time) (err error)
}

package statemanager

import (
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type stateManager struct {
	stateService interfaces.StateService
}

func NewStateManager(stateService interfaces.StateService) interfaces.StateManager {
	return &stateManager{
		stateService: stateService,
	}
}

func (ref *stateManager) GetOrCreateState(eventDate time.Time, entityId uuid.UUID, engineName values.EngineName) (state *entity.State, shouldIgnore bool, err error) {
	state, exists, err := ref.stateService.Get(entityId)
	if err != nil {
		return nil, false, errors.WithStack(err)
	}
	if exists && eventDate.Before(state.RequestDate) {
		logrus.WithField("entityId", entityId).
			Infof("[commandAndControl.getOrCreateState] Ignoring the execution because there is a "+
				"more recent engine execution, requestDate: %s - last requestDate: %s", eventDate.String(), state.RequestDate.String())

		return state, true, nil
	}
	if !exists {
		state, err = ref.stateService.Create(entityId, engineName)
		if err != nil {
			return nil, false, errors.WithStack(err)
		}
	}

	return state, false, nil
}

func (ref *stateManager) UpdateState(state *entity.State, requestDate time.Time, executionTime time.Time) (err error) {
	return ref.stateService.Update(*state, requestDate, executionTime)
}

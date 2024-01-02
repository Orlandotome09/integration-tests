package stateEventPublisher

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventPublisher/stateEventPublisher/contract"
)

type stateEventPublisher struct {
	profileRepository  interfaces.ComplianceProfileRepository
	contractRepository interfaces.ContractRepository
	personRepository   interfaces.PersonRepository
	queuePublisher     interfaces.QueuePublisher
	idGenerator        func() uuid.UUID
}

func NewStateEventsPublisher(
	profileRepository interfaces.ComplianceProfileRepository,
	contractRepository interfaces.ContractRepository,
	personRepository interfaces.PersonRepository,
	queuePublisher interfaces.QueuePublisher,
	idGenerator func() uuid.UUID,
) interfaces.StateEventsPublisher {
	return &stateEventPublisher{
		profileRepository:  profileRepository,
		contractRepository: contractRepository,
		personRepository:   personRepository,
		queuePublisher:     queuePublisher,
		idGenerator:        idGenerator,
	}
}

func (ref *stateEventPublisher) Send(state entity.State, eventType values.EventType) error {
	var stateEvent contract.StateEvent

	switch state.EngineName {
	case values.EngineNameProfile:
		profile, err := ref.profileRepository.Get(state.EntityID)
		if err != nil {
			return errors.WithStack(err)
		}
		stateEvent = contract.NewProfileStateEvent(state, profile, eventType)
	case values.EngineNamePerson:
		person, err := ref.personRepository.Get(state.EntityID)
		if err != nil {
			return errors.WithStack(err)
		}
		stateEvent = contract.NewPersonStateEvent(state, person, eventType)
	case values.EngineNameContract:
		contractRow, err := ref.contractRepository.Get(state.EntityID)
		if err != nil {
			return errors.WithStack(err)
		}
		stateEvent = contract.NewContractStateEvent(state, contractRow, eventType)
	default:
		logrus.Infof("[StateEventsPublisher] not publishing event for %s. "+
			"Only publishing events for PROFILE, PERSON and CONTRACT", state.EngineName)

		return nil
	}

	stateEvent.ID = ref.idGenerator().String()
	message, err := json.Marshal(stateEvent)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := ref.queuePublisher.Publish(message, "", stateEvent.EntityID); err != nil {
		logrus.WithField("state_event", string(message)).
			Errorf("[StateEventsPublisher] error publishing state event with entity id %s", stateEvent.EntityID)
		return errors.WithStack(err)
	}

	logrus.WithField("state_event", string(message)).
		Infof("[StateEventsPublisher] published state event with entity id %s", stateEvent.EntityID)

	return nil
}

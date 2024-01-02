package stateEventPublisher

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/eventPublisher/stateEventPublisher/contract"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	id          = uuid.New()
	idGenerator = func() uuid.UUID { return id }
)

type params struct {
	profileRepository  *mocks.ComplianceProfileRepository
	contractRepository *mocks.ContractRepository
	personRepository   *mocks.PersonRepository
	queuePublisher     *mocks.QueuePublisher
	idGenerator        func() uuid.UUID
}

func prepare() (params, interfaces.StateEventsPublisher) {
	profileRepository := &mocks.ComplianceProfileRepository{}
	contractRepository := &mocks.ContractRepository{}
	personRepository := &mocks.PersonRepository{}
	queuePublisher := &mocks.QueuePublisher{}
	publisher := NewStateEventsPublisher(profileRepository, contractRepository, personRepository, queuePublisher, idGenerator)
	return params{
		profileRepository:  profileRepository,
		contractRepository: contractRepository,
		personRepository:   personRepository,
		queuePublisher:     queuePublisher,
		idGenerator:        idGenerator,
	}, publisher
}

func TestSend_ShouldSendProfileEvent(t *testing.T) {
	params, publisher := prepare()

	state := entity.State{
		EntityID:   uuid.New(),
		EngineName: values.EngineNameProfile,
	}
	eventType := values.EventTypeStateCreated
	profileID := uuid.New()
	profile := &entity.Profile{ProfileID: &profileID}

	stateEvent := contract.NewProfileStateEvent(state, profile, eventType)
	stateEvent.ID = idGenerator().String()
	message, _ := json.Marshal(stateEvent)

	params.profileRepository.On("Get", state.EntityID).Return(profile, nil)
	params.queuePublisher.On("Publish", message, "", stateEvent.EntityID).Return(nil)

	err := publisher.Send(state, eventType)

	assert.Nil(t, err)
}

func TestSend_ShouldSendPersonEvent(t *testing.T) {
	params, publisher := prepare()

	state := entity.State{
		EntityID:   uuid.New(),
		EngineName: values.EngineNamePerson,
	}
	eventType := values.EventTypeStateCreated
	person := &entity.Person{EntityID: uuid.New()}

	stateEvent := contract.NewPersonStateEvent(state, person, eventType)
	stateEvent.ID = idGenerator().String()
	message, _ := json.Marshal(stateEvent)

	params.personRepository.On("Get", state.EntityID).Return(person, nil)
	params.queuePublisher.On("Publish", message, "", stateEvent.EntityID).Return(nil)

	err := publisher.Send(state, eventType)

	assert.Nil(t, err)
}

func TestSend_ShouldSendContractEvent(t *testing.T) {
	params, publisher := prepare()
	profileID := uuid.New()
	state := entity.State{
		EntityID:   uuid.New(),
		EngineName: values.EngineNameContract,
	}
	eventType := values.EventTypeStateCreated
	row := &entity.Contract{
		ContractID:           &state.EntityID,
		EstimatedTotalAmount: 0,
		DueTime:              "",
		Installments:         0,
		CorrelationID:        "",
		ProfileID:            &profileID,
		DocumentID:           nil,
	}

	stateEvent := contract.NewContractStateEvent(state, row, eventType)
	stateEvent.ID = idGenerator().String()
	message, _ := json.Marshal(stateEvent)

	params.contractRepository.On("Get", state.EntityID).Return(row, nil)
	params.queuePublisher.On("Publish", message, "", stateEvent.EntityID).Return(nil)

	err := publisher.Send(state, eventType)

	assert.Nil(t, err)
}

func TestSend_ShouldNotSendEventForContractEngine(t *testing.T) {
	params, publisher := prepare()
	profileID := uuid.New()

	state := entity.State{
		EntityID:   uuid.New(),
		EngineName: values.EngineNameContract,
	}
	eventType := values.EventTypeStateCreated

	row := &entity.Contract{
		ContractID:           &state.EntityID,
		EstimatedTotalAmount: 0,
		DueTime:              "",
		Installments:         0,
		CorrelationID:        "",
		ProfileID:            &profileID,
		DocumentID:           nil,
	}

	stateEvent := contract.NewContractStateEvent(state, row, eventType)
	stateEvent.ID = idGenerator().String()
	message, _ := json.Marshal(stateEvent)

	params.contractRepository.On("Get", state.EntityID).Return(row, nil)
	params.queuePublisher.On("Publish", message, "", stateEvent.EntityID).Return(nil)

	err := publisher.Send(state, eventType)

	assert.Nil(t, err)
	params.profileRepository.AssertNumberOfCalls(t, "Get", 0)
	params.personRepository.AssertNumberOfCalls(t, "Get", 0)
	params.contractRepository.AssertNumberOfCalls(t, "Get", 1)
}

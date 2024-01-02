package state

import (
	"testing"
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type params struct {
	stateRepository      *mocks.StateRepository
	profileAdapter       *mocks.ProfileAdapter
	contractAdapter      *mocks.ContractAdapter
	stateEventsPublisher *mocks.StateEventsPublisher
}

func prepare() (params, interfaces.StateService) {
	stateRepository := &mocks.StateRepository{}
	profileService := &mocks.ProfileAdapter{}
	contractAdapter := &mocks.ContractAdapter{}
	stateEventsPublisher := &mocks.StateEventsPublisher{}
	complianceCommandPublisher := &mocks.ComplianceCommandPublisher{}
	service := NewStateService(stateRepository, profileService, contractAdapter, stateEventsPublisher, complianceCommandPublisher)
	return params{
		stateRepository:      stateRepository,
		profileAdapter:       profileService,
		contractAdapter:      contractAdapter,
		stateEventsPublisher: stateEventsPublisher,
	}, service
}

func TestGet(t *testing.T) {
	params, service := prepare()

	entityID := uuid.New()
	state := &entity.State{}

	params.stateRepository.On("Get", entityID).Return(state, true, nil)

	receivedState, exists, err := service.Get(entityID)

	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, state, receivedState)
}

func TestCreate(t *testing.T) {
	params, service := prepare()

	entityID := uuid.New()
	engineName := "PROFILE"
	state := &entity.State{}

	params.stateRepository.On("Create", entityID, engineName).Return(state, nil)
	params.stateEventsPublisher.On("Send", *state, values.EventTypeStateCreated).Return(nil)

	receivedState, err := service.Create(entityID, engineName)

	assert.Nil(t, err)
	assert.Equal(t, state, receivedState)
}

func TestUpdate(t *testing.T) {
	params, service := prepare()

	state := entity.State{}
	saved := &entity.State{}
	requestDate := time.Now()
	executionTime := time.Now()

	params.stateRepository.On("Save", state, requestDate, executionTime).Return(saved, nil)
	params.stateEventsPublisher.On("Send", *saved, values.EventTypeStateChanged).Return(nil)

	err := service.Update(state, requestDate, executionTime)

	assert.Nil(t, err)
}

func TestUpdate_shouldDiscardEvent(t *testing.T) {
	params, service := prepare()

	state := entity.State{}
	var saved *entity.State
	requestDate := time.Now()
	executionTime := time.Now()

	params.stateRepository.On("Save", state, requestDate, executionTime).Return(saved, nil)

	err := service.Update(state, requestDate, executionTime)

	assert.Nil(t, err)
	params.stateEventsPublisher.AssertNumberOfCalls(t, "Send", 0)
}

func TestSearchProfileStates(t *testing.T) {
	params, service := prepare()

	request := entity.SearchProfileStateRequest{}
	states := []entity.State{{EntityID: uuid.New()}, {EntityID: uuid.New()}}
	parentID0 := uuid.New()
	profileID0 := uuid.New()
	profile0 := &entity.Profile{
		ProfileID: &profileID0,
		Person: entity.Person{
			DocumentNumber: "111",
			Name:           "First",
			PersonType:     values.PersonTypeIndividual,
			PartnerID:      "111",
			Individual:     &entity.Individual{FirstName: "First", LastName: "Last"},
		},
		ParentID: &parentID0,
	}
	profileID1 := uuid.New()
	parentID1 := uuid.New()
	profile1 := &entity.Profile{
		ProfileID: &profileID1,
		ParentID:  &parentID1,
		Person: entity.Person{
			DocumentNumber: "222",
			Name:           "Second",
			PersonType:     values.PersonTypeCompany,
			PartnerID:      "222",
		},
	}

	params.stateRepository.On("Search", request).Return(states, int64(2), nil)
	params.profileAdapter.On("Get", states[0].EntityID).Return(profile0, nil)
	params.profileAdapter.On("Get", states[1].EntityID).Return(profile1, nil)

	expected := &entity.ProfileStateList{
		Count: 2,
		ProfileStates: []entity.ProfileState{
			{State: states[0], ProfileID: *profile0.ProfileID, PartnerID: profile0.PartnerID, ParentID: profile0.ParentID, DocumentNumber: profile0.DocumentNumber, Name: profile0.Name},
			{State: states[1], ProfileID: *profile1.ProfileID, PartnerID: profile1.PartnerID, ParentID: profile1.ParentID, DocumentNumber: profile1.DocumentNumber, Name: profile1.Name},
		},
	}

	received, err := service.SearchProfileStates(request)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestSearchContractStates(t *testing.T) {
	params, service := prepare()

	request := entity.SearchProfileStateRequest{}
	states := []entity.State{{EntityID: uuid.New()}, {EntityID: uuid.New()}}
	contractID0 := uuid.New()
	contract0 := &entity.Contract{ProfileID: &contractID0}
	contractID1 := uuid.New()
	contract1 := &entity.Contract{ProfileID: &contractID1}
	profileID0 := uuid.New()
	parentID0 := uuid.New()

	profile0 := &entity.Profile{
		ProfileID: &profileID0,
		ParentID:  &parentID0,
		Person: entity.Person{
			DocumentNumber: "doc0",
			Name:           "name0",
			PersonType:     values.PersonTypeCompany,
			PartnerID:      "partner0",
		},
	}

	profileID1 := uuid.New()
	parentID1 := uuid.New()
	profile1 := &entity.Profile{
		ProfileID: &profileID1,
		ParentID:  &parentID1,
		Person: entity.Person{
			Name:           "Peterson",
			DocumentNumber: "doc1",
			PartnerID:      "partner1",
			PersonType:     values.PersonTypeIndividual,
			Individual: &entity.Individual{
				FirstName: "first",
				LastName:  "last",
			},
		},
	}

	params.stateRepository.On("SearchContractStates", request).Return(states, int64(2), nil)
	params.contractAdapter.On("Get", &states[0].EntityID).Return(contract0, true, nil)
	params.contractAdapter.On("Get", &states[1].EntityID).Return(contract1, true, nil)
	params.profileAdapter.On("Get", *contract0.ProfileID).Return(profile0, nil)
	params.profileAdapter.On("Get", *contract1.ProfileID).Return(profile1, nil)

	expected := &entity.ProfileStateList{
		Count: 2,
		ProfileStates: []entity.ProfileState{
			{State: states[0], ProfileID: *profile0.ProfileID, PartnerID: profile0.PartnerID, ParentID: profile0.ParentID, DocumentNumber: profile0.DocumentNumber, Name: profile0.Name},
			{State: states[1], ProfileID: *profile1.ProfileID, PartnerID: profile1.PartnerID, ParentID: profile1.ParentID, DocumentNumber: profile1.DocumentNumber, Name: profile1.Name},
		},
	}

	received, err := service.SearchContractStates(request)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestFindByProfileID(t *testing.T) {
	params, service := prepare()

	profileID := uuid.New()
	engine := values.EngineNameProfile
	result := values.ResultStatusAnalysing
	states := []entity.State{}

	params.stateRepository.On("FindByProfileID", profileID, engine, result).Return(states, nil)

	expected := states
	received, err := service.FindByProfileID(profileID, engine, result)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestResync(t *testing.T) {
	params, service := prepare()
	t.Run("no ids to resync", func(t *testing.T) {
		resynced, err := service.Resync()

		a := assert.New(t)

		a.Nil(resynced)
		a.Nil(err)
	})

	t.Run("resync", func(t *testing.T) {
		a := assert.New(t)

		entityIds := []string{}
		for i := 0; i < 100; i++ {
			entityIds = append(entityIds, uuid.New().String())
		}

		for _, entityId := range entityIds {
			state := &entity.State{
				EntityID:   uuid.MustParse(entityId),
				EngineName: values.EngineNameProfile,
			}
			params.stateRepository.On("Get", uuid.MustParse(entityId)).Return(state, true, nil)
			params.stateEventsPublisher.On("Send", *state, values.EventTypeStateResync).
				Return(nil)
		}

		resynced, err := service.Resync(entityIds...)

		a.Nil(err)
		a.Equal(len(entityIds), len(resynced))

	})

}

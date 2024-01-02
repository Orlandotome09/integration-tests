package contract

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewProfileStateEvent(t *testing.T) {
	profileID := uuid.New()
	state := entity.State{
		EntityID:               profileID,
		EngineName:             values.EngineNameProfile,
		Result:                 values.ResultStatusApproved,
		ValidationStepsResults: nil,
		RuleNames:              nil,
		Pending:                false,
		ExecutionTime:          time.Now(),
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	parentID := uuid.New()
	profile := entity.Profile{
		Person:               entity.Person{DocumentNumber: uuid.New().String()},
		ProfileID:            &profileID,
		ParentID:             &parentID,
		LegacyID:             uuid.New().String(),
		CallbackUrl:          "/callback",
		LegalRepresentatives: nil,
		OwnershipStructure:   nil,
		BoardOfDirectors:     nil,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	eventType := values.EventTypeStateCreated

	expected := StateEvent{
		EventHeader: entity.EventHeader{
			EntityID:   state.EntityID.String(),
			EntityType: values.EntityTypeComplianceState,
			EventType:  eventType,
			UpdateDate: state.UpdatedAt,
		},
		Data: StateEventData{
			ProfileID: *profile.ProfileID,
			Content: StateEventContent{
				State:   state,
				Profile: &profile,
			},
		},
	}

	received := NewProfileStateEvent(state, &profile, eventType)

	assert.Equal(t, expected, received)
}

func TestNewProfileStateEventForProfileNil(t *testing.T) {
	profileID := uuid.New()
	state := entity.State{
		EntityID:               profileID,
		EngineName:             values.EngineNameProfile,
		Result:                 values.ResultStatusApproved,
		ValidationStepsResults: nil,
		RuleNames:              nil,
		Pending:                false,
		ExecutionTime:          time.Now(),
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	var profile *entity.Profile = nil
	eventType := values.EventTypeStateCreated

	expected := StateEvent{
		EventHeader: entity.EventHeader{
			EntityID:   state.EntityID.String(),
			EntityType: values.EntityTypeComplianceState,
			EventType:  eventType,
			UpdateDate: state.UpdatedAt,
		},
		Data: StateEventData{
			ProfileID: state.EntityID,
			Content: StateEventContent{
				State:   state,
				Profile: profile,
			},
		},
	}

	received := NewProfileStateEvent(state, profile, eventType)

	assert.Equal(t, expected, received)
}

func TestNewPersonStateEvent(t *testing.T) {
	personID := uuid.New()
	state := entity.State{
		EntityID:               personID,
		EngineName:             values.EngineNameProfile,
		Result:                 values.ResultStatusApproved,
		ValidationStepsResults: nil,
		RuleNames:              nil,
		Pending:                false,
		ExecutionTime:          time.Now(),
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	person := entity.Person{
		DocumentNumber: uuid.New().String(),
		Name:           uuid.New().String(),
		PersonType:     values.PersonTypeIndividual,
		Email:          uuid.New().String(),
		PartnerID:      uuid.New().String(),
		OfferType:      uuid.New().String(),
		ProfileID:      uuid.New(),
		EntityID:       personID,
		EntityType:     values.EntityTypeDirector,
		RoleType:       values.RoleTypeDirector,
	}
	eventType := values.EventTypeStateCreated

	expected := StateEvent{
		EventHeader: entity.EventHeader{
			EntityID:   state.EntityID.String(),
			EntityType: values.EntityTypeComplianceState,
			EventType:  eventType,
			UpdateDate: state.UpdatedAt,
		},
		Data: StateEventData{
			ProfileID: person.ProfileID,
			Content: StateEventContent{
				State:  state,
				Person: &person,
			},
		},
	}

	received := NewPersonStateEvent(state, &person, eventType)

	assert.Equal(t, expected, received)
}

func TestNewPersonStateEventForPersonNil(t *testing.T) {
	personID := uuid.New()
	state := entity.State{
		EntityID:               personID,
		EngineName:             values.EngineNameProfile,
		Result:                 values.ResultStatusApproved,
		ValidationStepsResults: nil,
		RuleNames:              nil,
		Pending:                false,
		ExecutionTime:          time.Now(),
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	var person *entity.Person = nil
	eventType := values.EventTypeStateCreated

	expected := StateEvent{
		EventHeader: entity.EventHeader{
			EntityID:   state.EntityID.String(),
			EntityType: values.EntityTypeComplianceState,
			EventType:  eventType,
			UpdateDate: state.UpdatedAt,
		},
		Data: StateEventData{
			ProfileID: uuid.Nil,
			Content: StateEventContent{
				State:  state,
				Person: person,
			},
		},
	}

	received := NewPersonStateEvent(state, person, eventType)

	assert.Equal(t, expected, received)
}

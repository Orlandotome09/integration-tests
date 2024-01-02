package contract

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type StateEvent struct {
	entity.EventHeader
	Data StateEventData `json:"data"`
}

type StateEventData struct {
	ProfileID uuid.UUID         `json:"profile_id"`
	Content   StateEventContent `json:"content"`
}

type StateEventContent struct {
	State    entity.State     `json:"state"`
	Profile  *entity.Profile  `json:"profile"`
	Person   *entity.Person   `json:"person"`
	Contract *entity.Contract `json:"contract"`
}

func NewProfileStateEvent(state entity.State, profile *entity.Profile, eventType values.EventType) StateEvent {
	return StateEvent{
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
}

func NewPersonStateEvent(state entity.State, person *entity.Person, eventType values.EventType) StateEvent {
	var profileID uuid.UUID
	if person != nil {
		profileID = person.ProfileID
	}

	return StateEvent{
		EventHeader: entity.EventHeader{
			EntityID:   state.EntityID.String(),
			EntityType: values.EntityTypeComplianceState,
			EventType:  eventType,
			UpdateDate: state.UpdatedAt,
		},
		Data: StateEventData{
			ProfileID: profileID,
			Content: StateEventContent{
				State:  state,
				Person: person,
			},
		},
	}
}

func NewContractStateEvent(state entity.State, contract *entity.Contract, eventType values.EventType) StateEvent {
	var profileID uuid.UUID
	if contract != nil {
		profileID = *contract.ProfileID
	}
	return StateEvent{
		EventHeader: entity.EventHeader{
			EntityID:   state.EntityID.String(),
			EntityType: values.EntityTypeComplianceState,
			EventType:  eventType,
			UpdateDate: state.UpdatedAt,
		},
		Data: StateEventData{
			ProfileID: profileID,
			Content: StateEventContent{
				State:    state,
				Contract: contract,
			},
		},
	}
}

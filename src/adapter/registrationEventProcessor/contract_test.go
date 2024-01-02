package registrationEventProcessor

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationEvent_IsValid(t *testing.T) {
	a := assert.New(t)

	type fields struct {
		EventType  string
		ProfileID  uuid.UUID
		EntityID   uuid.UUID
		EntityType string
		ParentID   *uuid.UUID
		ParentType string
		UpdateDate time.Time
		Content    json.RawMessage
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "Should be valid when UpdateDate is not zero and EntityID is not nil",
			fields: fields{
				EventType:  "PROFILE_CHANGED",
				ProfileID:  uuid.New(),
				EntityID:   uuid.New(),
				EntityType: "PROFILE",
				ParentID:   nil,
				ParentType: "PROFILE",
				UpdateDate: time.Now(),
				Content:    nil,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should be valid when UpdateDate is not zero and EntityID is nil",
			fields: fields{
				EventType:  "PROFILE_CHANGED",
				ProfileID:  uuid.New(),
				EntityID:   uuid.Nil,
				EntityType: "PROFILE",
				ParentID:   nil,
				ParentType: "PROFILE",
				UpdateDate: time.Now(),
				Content:    nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Should be valid when UpdateDate is zero and EntityID is not nil",
			fields: fields{
				EventType:  "ADDRESS_CREATED",
				ProfileID:  uuid.New(),
				EntityID:   uuid.New(),
				EntityType: "ADDRESS",
				ParentID:   nil,
				ParentType: "PROFILE",
				UpdateDate: time.Time{},
				Content:    nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := RegistrationEvent{
				EventType:  tt.fields.EventType,
				ProfileID:  tt.fields.ProfileID,
				EntityID:   tt.fields.EntityID,
				EntityType: tt.fields.EntityType,
				ParentID:   tt.fields.ParentID,
				ParentType: tt.fields.ParentType,
				UpdateDate: tt.fields.UpdateDate,
				Content:    tt.fields.Content,
			}
			got, err := event.IsValid()

			a.Equal(tt.wantErr, err != nil)
			a.Equal(tt.want, got)
		})
	}
}

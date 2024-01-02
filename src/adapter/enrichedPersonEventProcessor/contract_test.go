package enrichedPersonEventProcessor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrichedPersonEvent_IsValid(t *testing.T) {
	a := assert.New(t)

	type fields struct {
		EntityID   uuid.UUID
		EntityType string
		EventType  string
		Data       PersonEnrichedData
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "should return false when event without profile id",
			fields: fields{
				EntityID: uuid.Nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "should return true when event with profile id",
			fields: fields{
				EntityID: uuid.New(),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := EnrichedPersonEvent{
				EntityID:   tt.fields.EntityID,
				EntityType: tt.fields.EntityType,
				EventType:  tt.fields.EventType,
				Data:       tt.fields.Data,
			}
			got, err := event.IsValid()

			a.Equal(tt.wantErr, err != nil)
			a.Equal(tt.want, got)
		})
	}
}

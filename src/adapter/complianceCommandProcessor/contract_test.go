package complianceCommandProcessor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrichedPersonEvent_IsValid(t *testing.T) {
	a := assert.New(t)
	type fields struct {
		EntityID  uuid.UUID
		EventType string
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
			event := complianceCommand{
				EntityID:  tt.fields.EntityID,
				EventType: tt.fields.EventType,
			}
			got, err := event.IsValid()

			a.Equal(tt.wantErr, err != nil)
			a.Equal(tt.want, got)
		})
	}
}

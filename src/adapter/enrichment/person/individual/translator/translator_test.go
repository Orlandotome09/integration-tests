package enrichedIndividualTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/enrichment/person/individual/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := New()

	response := &contracts.IndividualResponse{
		Name:      "Sabrina Yasmin Oliveira",
		BirthDate: "26/01/1976",
		Situation: 1,
	}

	responseBytes := new(bytes.Buffer)
	json.NewEncoder(responseBytes).Encode(response)

	expected := &entity.EnrichedInformation{
		BureauStatus: "REGULAR",
		EnrichedIndividual: entity.EnrichedIndividual{
			Name:      response.Name,
			BirthDate: response.BirthDate,
		},
	}

	received, err := translator.Translate("123", uuid.New(), responseBytes.Bytes())

	assert.Nil(t, err)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func Test_enrichedIndividualTranslator_translateSituationToStatus(t *testing.T) {
	type args struct {
		situation int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should translate deceased status",
			args: args{situation: 0},
			want: "DECEASED",
		},
		{
			name: "should translate regular status",
			args: args{situation: 1},
			want: "REGULAR",
		},
		{
			name: "should translate pending status",
			args: args{situation: 2},
			want: "PENDING_REGULARIZATION",
		},
		{
			name: "should translate suspended status",
			args: args{situation: 3},
			want: "SUSPENDED",
		},
		{
			name: "should translate canceled status",
			args: args{situation: 4},
			want: "CANCELLED",
		},
		{
			name: "should translate null status",
			args: args{situation: 5},
			want: "NULL",
		},
		{
			name: "should translate unknown status",
			args: args{situation: 6},
			want: "UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := &enrichedIndividualTranslator{}
			assert.Equalf(t, tt.want, ref.translateSituationToStatus(tt.args.situation), "translateSituationToStatus(%v)", tt.args.situation)
		})
	}
}

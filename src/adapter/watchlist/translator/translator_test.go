package watchlistTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http/dto"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"
	"reflect"
	"testing"
)

func TestToDomain(t *testing.T) {
	translator := New()

	responses := []dto.WatchlistResponse{
		{Sources: []string{"PEP", "OFAC"}}, {Sources: []string{"TERROR", "BACEN"}},
	}
	metadata, _ := json.Marshal(responses)

	received, err := translator.ToDomain(responses)
	expected := &entity2.Watchlist{Sources: []string{"PEP", "OFAC", "TERROR", "BACEN"}, Metadata: metadata}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", nil, err)
	}
}

func TestToDomainMixedCase(t *testing.T) {
	translator := New()

	responses := []dto.WatchlistResponse{
		{Sources: []string{"Enforcement"}}, {Sources: []string{"Ofac"}},
	}
	metadata, _ := json.Marshal(responses)

	received, err := translator.ToDomain(responses)
	expected := &entity2.Watchlist{Sources: []string{"ENFORCEMENT", "OFAC"}, Metadata: metadata}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", nil, err)
	}
}

func TestToDomain_emptySources(t *testing.T) {
	translator := New()

	responses := []dto.WatchlistResponse{}
	metadata, _ := json.Marshal(responses)

	received, err := translator.ToDomain(responses)
	expected := &entity2.Watchlist{Sources: []string{}, Metadata: metadata}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", nil, err)
	}
}

func Test_translateSource(t *testing.T) {
	type args struct {
		inputSource string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should translate to uppercase",
			args: args{inputSource: "someMixedCase"},
			want: "SOMEMIXEDCASE",
		},
		{
			name: "should translate spaces to underline",
			args: args{inputSource: "Adverse Media"},
			want: "ADVERSE_MEDIA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := translateSource(tt.args.inputSource); got != tt.want {
				t.Errorf("translateSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

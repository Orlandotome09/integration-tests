package blacklistTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http/dto"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"reflect"
	"testing"
	"time"
)

func TestToDomain(t *testing.T) {
	translator := New()

	timeNow := time.Now()
	response := &dto.BlacklistResponse{
		CreatedAt:      timeNow,
		Status:         "ACTIVE",
		Justifications: []dto.Justification{{Justification: "AAA"}, {Justification: "BBB"}},
	}

	received := translator.ToDomain(response)

	expected := &entity.BlacklistStatus{
		Status: response.Status,
		Justification: entity.Justification{
			AddedAt:  timeNow,
			Author:   "",
			Comments: []string{"AAA", "BBB"},
		},
	}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

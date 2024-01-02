package timeTranslator

import (
	"reflect"
	"testing"
	"time"
)

func TestTranslateTime(t *testing.T) {
	translator := NewTimeTranslator()

	date := time.Date(2024, 01, 01, 11, 1, 1, 1, time.UTC)

	expected := "01012024"

	received := translator.TranslateTime(date)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

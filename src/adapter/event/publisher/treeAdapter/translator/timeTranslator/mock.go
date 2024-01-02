package timeTranslator

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MockTimeTranslator struct {
	TimeTranslator
	mock.Mock
}

func (ref *MockTimeTranslator) TranslateTime(dateTime time.Time) string {
	args := ref.Called(dateTime)
	return args.Get(0).(string)
}

func (ref *MockTimeTranslator) GenerateMessageDate() time.Time {
	args := ref.Called()
	return args.Get(0).(time.Time)
}

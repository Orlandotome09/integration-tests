package message

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type Mock struct {
	mock.Mock
}

func (ref *Mock) ID() string {
	ret := ref.Called()
	return ret.Get(0).(string)
}

func (ref *Mock) Data() []byte {
	ret := ref.Called()
	return ret.Get(0).([]byte)
}

func (ref *Mock) PublishTime() time.Time {
	ret := ref.Called()
	return ret.Get(0).(time.Time)
}

func (ref *Mock) Ack(result string) {
	ref.Called(result)
	return
}

func (ref *Mock) Nack(result string) bool {
	ret := ref.Called(result)
	return ret.Get(0).(bool)
}

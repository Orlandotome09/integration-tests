// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// QueuePublisher is an autogenerated mock type for the QueuePublisher type
type QueuePublisher struct {
	mock.Mock
}

// Publish provides a mock function with given fields: message, orderingKey, concurrencyKey
func (_m *QueuePublisher) Publish(message []byte, orderingKey string, concurrencyKey string) error {
	ret := _m.Called(message, orderingKey, concurrencyKey)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, string, string) error); ok {
		r0 = rf(message, orderingKey, concurrencyKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

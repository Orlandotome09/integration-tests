// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import _interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
import mock "github.com/stretchr/testify/mock"

// QueueSubscriber is an autogenerated mock type for the QueueSubscriber type
type QueueSubscriber struct {
	mock.Mock
}

// Subscribe provides a mock function with given fields: processor
func (_m *QueueSubscriber) Subscribe(processor _interfaces.MessageProcessor) error {
	ret := _m.Called(processor)

	var r0 error
	if rf, ok := ret.Get(0).(func(_interfaces.MessageProcessor) error); ok {
		r0 = rf(processor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscriptionName provides a mock function with given fields:
func (_m *QueueSubscriber) SubscriptionName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

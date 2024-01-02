// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import _interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
import mock "github.com/stretchr/testify/mock"

// EventListener is an autogenerated mock type for the EventListener type
type EventListener struct {
	mock.Mock
}

// Listen provides a mock function with given fields:
func (_m *EventListener) Listen() {
	_m.Called()
}

// Register provides a mock function with given fields: subscriber, processor
func (_m *EventListener) Register(subscriber _interfaces.QueueSubscriber, processor _interfaces.MessageProcessor) {
	_m.Called(subscriber, processor)
}
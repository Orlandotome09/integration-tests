// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"

// StateEventsPublisher is an autogenerated mock type for the StateEventsPublisher type
type StateEventsPublisher struct {
	mock.Mock
}

// Send provides a mock function with given fields: state, eventType
func (_m *StateEventsPublisher) Send(state entity.State, eventType string) error {
	ret := _m.Called(state, eventType)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.State, string) error); ok {
		r0 = rf(state, eventType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

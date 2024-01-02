// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import _interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
import mock "github.com/stretchr/testify/mock"

// Subscriber is an autogenerated mock type for the Subscriber type
type Subscriber struct {
	mock.Mock
}

// Listen provides a mock function with given fields: processor
func (_m *Subscriber) Listen(processor _interfaces.Processor) error {
	ret := _m.Called(processor)

	var r0 error
	if rf, ok := ret.Get(0).(func(_interfaces.Processor) error); ok {
		r0 = rf(processor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
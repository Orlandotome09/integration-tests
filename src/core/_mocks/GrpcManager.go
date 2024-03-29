// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import _interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
import mock "github.com/stretchr/testify/mock"

// GrpcManager is an autogenerated mock type for the GrpcManager type
type GrpcManager struct {
	mock.Mock
}

// Listen provides a mock function with given fields: port, cncInstance
func (_m *GrpcManager) Listen(port int, cncInstance _interfaces.EventProcessor) error {
	ret := _m.Called(port, cncInstance)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, _interfaces.EventProcessor) error); ok {
		r0 = rf(port, cncInstance)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

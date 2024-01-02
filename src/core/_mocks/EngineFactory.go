// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import _interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
import mock "github.com/stretchr/testify/mock"

// EngineFactory is an autogenerated mock type for the EngineFactory type
type EngineFactory struct {
	mock.Mock
}

// CreateEngine provides a mock function with given fields: engineName
func (_m *EngineFactory) CreateEngine(engineName string) (_interfaces.Engine, error) {
	ret := _m.Called(engineName)

	var r0 _interfaces.Engine
	if rf, ok := ret.Get(0).(func(string) _interfaces.Engine); ok {
		r0 = rf(engineName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(_interfaces.Engine)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(engineName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

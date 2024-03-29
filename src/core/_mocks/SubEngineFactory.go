// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import _interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
import mock "github.com/stretchr/testify/mock"

// SubEngineFactory is an autogenerated mock type for the SubEngineFactory type
type SubEngineFactory struct {
	mock.Mock
}

// CreateSubEngine provides a mock function with given fields: subEngineName
func (_m *SubEngineFactory) CreateSubEngine(subEngineName string) (_interfaces.SubEngine, error) {
	ret := _m.Called(subEngineName)

	var r0 _interfaces.SubEngine
	if rf, ok := ret.Get(0).(func(string) _interfaces.SubEngine); ok {
		r0 = rf(subEngineName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(_interfaces.SubEngine)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(subEngineName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

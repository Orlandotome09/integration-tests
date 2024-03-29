// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import _interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"
import uuid "github.com/google/uuid"

// Engine is an autogenerated mock type for the Engine type
type Engine struct {
	mock.Mock
}

// GetName provides a mock function with given fields:
func (_m *Engine) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewInstance provides a mock function with given fields:
func (_m *Engine) NewInstance() _interfaces.Engine {
	ret := _m.Called()

	var r0 _interfaces.Engine
	if rf, ok := ret.Get(0).(func() _interfaces.Engine); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(_interfaces.Engine)
		}
	}

	return r0
}

// PosProcessing provides a mock function with given fields: previousState, newState, entityID
func (_m *Engine) PosProcessing(previousState *entity.State, newState *entity.State, entityID uuid.UUID) error {
	ret := _m.Called(previousState, newState, entityID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.State, *entity.State, uuid.UUID) error); ok {
		r0 = rf(previousState, newState, entityID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Prepare provides a mock function with given fields: entityID
func (_m *Engine) Prepare(entityID uuid.UUID) error {
	ret := _m.Called(entityID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(entityID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Validate provides a mock function with given fields: state, override, noCache, entityID, engineName
func (_m *Engine) Validate(state entity.State, override entity.Overrides, noCache bool, entityID uuid.UUID, engineName string) (*entity.State, error) {
	ret := _m.Called(state, override, noCache, entityID, engineName)

	var r0 *entity.State
	if rf, ok := ret.Get(0).(func(entity.State, entity.Overrides, bool, uuid.UUID, string) *entity.State); ok {
		r0 = rf(state, override, noCache, entityID, engineName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.State, entity.Overrides, bool, uuid.UUID, string) error); ok {
		r1 = rf(state, override, noCache, entityID, engineName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

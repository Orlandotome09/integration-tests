// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"
import time "time"
import uuid "github.com/google/uuid"
import values "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"

// StateRepository is an autogenerated mock type for the StateRepository type
type StateRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: entityID, engineName
func (_m *StateRepository) Create(entityID uuid.UUID, engineName string) (*entity.State, error) {
	ret := _m.Called(entityID, engineName)

	var r0 *entity.State
	if rf, ok := ret.Get(0).(func(uuid.UUID, string) *entity.State); ok {
		r0 = rf(entityID, engineName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string) error); ok {
		r1 = rf(entityID, engineName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByProfileID provides a mock function with given fields: profileID, engine, result
func (_m *StateRepository) FindByProfileID(profileID uuid.UUID, engine string, result values.Result) ([]entity.State, error) {
	ret := _m.Called(profileID, engine, result)

	var r0 []entity.State
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, values.Result) []entity.State); ok {
		r0 = rf(profileID, engine, result)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string, values.Result) error); ok {
		r1 = rf(profileID, engine, result)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: entityID
func (_m *StateRepository) Get(entityID uuid.UUID) (*entity.State, bool, error) {
	ret := _m.Called(entityID)

	var r0 *entity.State
	if rf, ok := ret.Get(0).(func(uuid.UUID) *entity.State); ok {
		r0 = rf(entityID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.State)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(uuid.UUID) bool); ok {
		r1 = rf(entityID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(uuid.UUID) error); ok {
		r2 = rf(entityID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Save provides a mock function with given fields: state, requestDate, executionTime
func (_m *StateRepository) Save(state entity.State, requestDate time.Time, executionTime time.Time) (*entity.State, error) {
	ret := _m.Called(state, requestDate, executionTime)

	var r0 *entity.State
	if rf, ok := ret.Get(0).(func(entity.State, time.Time, time.Time) *entity.State); ok {
		r0 = rf(state, requestDate, executionTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.State, time.Time, time.Time) error); ok {
		r1 = rf(state, requestDate, executionTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: request
func (_m *StateRepository) Search(request entity.SearchProfileStateRequest) ([]entity.State, int64, error) {
	ret := _m.Called(request)

	var r0 []entity.State
	if rf, ok := ret.Get(0).(func(entity.SearchProfileStateRequest) []entity.State); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.State)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(entity.SearchProfileStateRequest) int64); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(entity.SearchProfileStateRequest) error); ok {
		r2 = rf(request)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SearchContractStates provides a mock function with given fields: request
func (_m *StateRepository) SearchContractStates(request entity.SearchProfileStateRequest) ([]entity.State, int64, error) {
	ret := _m.Called(request)

	var r0 []entity.State
	if rf, ok := ret.Get(0).(func(entity.SearchProfileStateRequest) []entity.State); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.State)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(entity.SearchProfileStateRequest) int64); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(entity.SearchProfileStateRequest) error); ok {
		r2 = rf(request)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"
import uuid "github.com/google/uuid"

// PosProcessor is an autogenerated mock type for the PosProcessor type
type PosProcessor struct {
	mock.Mock
}

// CreateInternalAccount provides a mock function with given fields: entityID
func (_m *PosProcessor) CreateInternalAccount(entityID uuid.UUID) error {
	ret := _m.Called(entityID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(entityID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendToLimit provides a mock function with given fields: profile, state, catalogConfig
func (_m *PosProcessor) SendToLimit(profile entity.Profile, state entity.State, catalogConfig entity.ProductConfig) error {
	ret := _m.Called(profile, state, catalogConfig)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.Profile, entity.State, entity.ProductConfig) error); ok {
		r0 = rf(profile, state, catalogConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendToTreeAdapter provides a mock function with given fields: profile, catalogConfig
func (_m *PosProcessor) SendToTreeAdapter(profile entity.Profile, catalogConfig entity.ProductConfig) error {
	ret := _m.Called(profile, catalogConfig)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.Profile, entity.ProductConfig) error); ok {
		r0 = rf(profile, catalogConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

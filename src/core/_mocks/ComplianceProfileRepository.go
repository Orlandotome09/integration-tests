// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"
import uuid "github.com/google/uuid"

// ComplianceProfileRepository is an autogenerated mock type for the ComplianceProfileRepository type
type ComplianceProfileRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: profileID
func (_m *ComplianceProfileRepository) Get(profileID uuid.UUID) (*entity.Profile, error) {
	ret := _m.Called(profileID)

	var r0 *entity.Profile
	if rf, ok := ret.Get(0).(func(uuid.UUID) *entity.Profile); ok {
		r0 = rf(profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Profile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(profileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: profile
func (_m *ComplianceProfileRepository) Save(profile entity.Profile) (*entity.Profile, error) {
	ret := _m.Called(profile)

	var r0 *entity.Profile
	if rf, ok := ret.Get(0).(func(entity.Profile) *entity.Profile); ok {
		r0 = rf(profile)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Profile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.Profile) error); ok {
		r1 = rf(profile)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

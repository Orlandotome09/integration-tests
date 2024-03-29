// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"
import uuid "github.com/google/uuid"

// LegalRepresentativeAdapter is an autogenerated mock type for the LegalRepresentativeAdapter type
type LegalRepresentativeAdapter struct {
	mock.Mock
}

// Get provides a mock function with given fields: id
func (_m *LegalRepresentativeAdapter) Get(id uuid.UUID) (*entity.LegalRepresentative, error) {
	ret := _m.Called(id)

	var r0 *entity.LegalRepresentative
	if rf, ok := ret.Get(0).(func(uuid.UUID) *entity.LegalRepresentative); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.LegalRepresentative)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: profileID
func (_m *LegalRepresentativeAdapter) Search(profileID uuid.UUID) ([]entity.LegalRepresentative, error) {
	ret := _m.Called(profileID)

	var r0 []entity.LegalRepresentative
	if rf, ok := ret.Get(0).(func(uuid.UUID) []entity.LegalRepresentative); ok {
		r0 = rf(profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.LegalRepresentative)
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

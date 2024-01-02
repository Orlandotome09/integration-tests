// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"

// AddressAdapter is an autogenerated mock type for the AddressAdapter type
type AddressAdapter struct {
	mock.Mock
}

// Get provides a mock function with given fields: id
func (_m *AddressAdapter) Get(id string) (*entity.Address, error) {
	ret := _m.Called(id)

	var r0 *entity.Address
	if rf, ok := ret.Get(0).(func(string) *entity.Address); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: profileID
func (_m *AddressAdapter) Search(profileID string) ([]entity.Address, error) {
	ret := _m.Called(profileID)

	var r0 []entity.Address
	if rf, ok := ret.Get(0).(func(string) []entity.Address); ok {
		r0 = rf(profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(profileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

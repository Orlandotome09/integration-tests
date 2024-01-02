// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"

// EconomicActivityService is an autogenerated mock type for the EconomicActivityService type
type EconomicActivityService struct {
	mock.Mock
}

// Get provides a mock function with given fields: code
func (_m *EconomicActivityService) Get(code string) (*entity.EconomicActivity, bool, error) {
	ret := _m.Called(code)

	var r0 *entity.EconomicActivity
	if rf, ok := ret.Get(0).(func(string) *entity.EconomicActivity); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.EconomicActivity)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(code)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
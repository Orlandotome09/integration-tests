// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"

// PartnerAdapter is an autogenerated mock type for the PartnerAdapter type
type PartnerAdapter struct {
	mock.Mock
}

// GetActive provides a mock function with given fields: partnerID
func (_m *PartnerAdapter) GetActive(partnerID string) (*entity.Partner, error) {
	ret := _m.Called(partnerID)

	var r0 *entity.Partner
	if rf, ok := ret.Get(0).(func(string) *entity.Partner); ok {
		r0 = rf(partnerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Partner)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(partnerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

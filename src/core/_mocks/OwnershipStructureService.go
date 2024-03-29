// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"

// OwnershipStructureService is an autogenerated mock type for the OwnershipStructureService type
type OwnershipStructureService struct {
	mock.Mock
}

// GetEnriched provides a mock function with given fields: legalEntityID, offerType, partnerID
func (_m *OwnershipStructureService) GetEnriched(legalEntityID string, offerType string, partnerID string) (*entity.OwnershipStructure, error) {
	ret := _m.Called(legalEntityID, offerType, partnerID)

	var r0 *entity.OwnershipStructure
	if rf, ok := ret.Get(0).(func(string, string, string) *entity.OwnershipStructure); ok {
		r0 = rf(legalEntityID, offerType, partnerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.OwnershipStructure)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(legalEntityID, offerType, partnerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManuallyFilled provides a mock function with given fields: profileID
func (_m *OwnershipStructureService) GetManuallyFilled(profileID string) (*entity.OwnershipStructure, error) {
	ret := _m.Called(profileID)

	var r0 *entity.OwnershipStructure
	if rf, ok := ret.Get(0).(func(string) *entity.OwnershipStructure); ok {
		r0 = rf(profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.OwnershipStructure)
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

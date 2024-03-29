// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"
import uuid "github.com/google/uuid"

// ContractAdapter is an autogenerated mock type for the ContractAdapter type
type ContractAdapter struct {
	mock.Mock
}

// Get provides a mock function with given fields: contractId
func (_m *ContractAdapter) Get(contractId *uuid.UUID) (*entity.Contract, bool, error) {
	ret := _m.Called(contractId)

	var r0 *entity.Contract
	if rf, ok := ret.Get(0).(func(*uuid.UUID) *entity.Contract); ok {
		r0 = rf(contractId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Contract)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(*uuid.UUID) bool); ok {
		r1 = rf(contractId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*uuid.UUID) error); ok {
		r2 = rf(contractId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

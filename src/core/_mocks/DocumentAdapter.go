// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"

// DocumentAdapter is an autogenerated mock type for the DocumentAdapter type
type DocumentAdapter struct {
	mock.Mock
}

// Find provides a mock function with given fields: entityID
func (_m *DocumentAdapter) Find(entityID string) ([]entity.Document, error) {
	ret := _m.Called(entityID)

	var r0 []entity.Document
	if rf, ok := ret.Get(0).(func(string) []entity.Document); ok {
		r0 = rf(entityID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Document)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(entityID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByEntityIDAndDocumentType provides a mock function with given fields: id, documentType
func (_m *DocumentAdapter) FindByEntityIDAndDocumentType(id string, documentType string) ([]entity.Document, error) {
	ret := _m.Called(id, documentType)

	var r0 []entity.Document
	if rf, ok := ret.Get(0).(func(string, string) []entity.Document); ok {
		r0 = rf(id, documentType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Document)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, documentType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *DocumentAdapter) GetByID(id string) (*entity.Document, error) {
	ret := _m.Called(id)

	var r0 *entity.Document
	if rf, ok := ret.Get(0).(func(string) *entity.Document); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Document)
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
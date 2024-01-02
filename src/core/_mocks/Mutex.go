// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Mutex is an autogenerated mock type for the Mutex type
type Mutex struct {
	mock.Mock
}

// Lock provides a mock function with given fields: id
func (_m *Mutex) Lock(id string) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Release provides a mock function with given fields: id
func (_m *Mutex) Release(id string) {
	_m.Called(id)
}

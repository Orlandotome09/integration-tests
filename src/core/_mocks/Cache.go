// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import time "time"

// Cache is an autogenerated mock type for the Cache type
type Cache struct {
	mock.Mock
}

// Get provides a mock function with given fields: key
func (_m *Cache) Get(key string) interface{} {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Save provides a mock function with given fields: key, value, ttl
func (_m *Cache) Save(key string, value interface{}, ttl time.Duration) {
	_m.Called(key, value, ttl)
}
// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"
import values "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"

// NotificationService is an autogenerated mock type for the NotificationService type
type NotificationService struct {
	mock.Mock
}

// SendNotification provides a mock function with given fields: event
func (_m *NotificationService) SendNotification(event *values.Event) (*entity.State, error) {
	ret := _m.Called(event)

	var r0 *entity.State
	if rf, ok := ret.Get(0).(func(*values.Event) *entity.State); ok {
		r0 = rf(event)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.State)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*values.Event) error); ok {
		r1 = rf(event)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

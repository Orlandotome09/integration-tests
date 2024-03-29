// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
import mock "github.com/stretchr/testify/mock"

// ProfileRulesFactory is an autogenerated mock type for the ProfileRulesFactory type
type ProfileRulesFactory struct {
	mock.Mock
}

// GetRules provides a mock function with given fields: ruleSetConfig, profile
func (_m *ProfileRulesFactory) GetRules(ruleSetConfig *entity.RuleSetConfig, profile *entity.Profile) []entity.Rule {
	ret := _m.Called(ruleSetConfig, profile)

	var r0 []entity.Rule
	if rf, ok := ret.Get(0).(func(*entity.RuleSetConfig, *entity.Profile) []entity.Rule); ok {
		r0 = rf(ruleSetConfig, profile)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Rule)
		}
	}

	return r0
}

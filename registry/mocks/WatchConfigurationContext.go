// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	client "github.com/coreos/etcd/client"
	mock "github.com/stretchr/testify/mock"
)

// WatchConfigurationContext is an autogenerated mock type for the WatchConfigurationContext type
type WatchConfigurationContext struct {
	mock.Mock
}

// Get provides a mock function with given fields: key
func (_m *WatchConfigurationContext) Get(key string) (string, error) {
	ret := _m.Called(key)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChildrenPaths provides a mock function with given fields: key
func (_m *WatchConfigurationContext) GetChildrenPaths(key string) ([]string, error) {
	ret := _m.Called(key)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Watch provides a mock function with given fields: key, recursive, eventChannel
func (_m *WatchConfigurationContext) Watch(key string, recursive bool, eventChannel chan *client.Response) {
	_m.Called(key, recursive, eventChannel)
}
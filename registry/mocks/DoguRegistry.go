// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	core "github.com/cloudogu/cesapp-lib/core"

	mock "github.com/stretchr/testify/mock"
)

// DoguRegistry is an autogenerated mock type for the DoguRegistry type
type DoguRegistry struct {
	mock.Mock
}

// Enable provides a mock function with given fields: dogu
func (_m *DoguRegistry) Enable(dogu *core.Dogu) error {
	ret := _m.Called(dogu)

	var r0 error
	if rf, ok := ret.Get(0).(func(*core.Dogu) error); ok {
		r0 = rf(dogu)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: name
func (_m *DoguRegistry) Get(name string) (*core.Dogu, error) {
	ret := _m.Called(name)

	var r0 *core.Dogu
	if rf, ok := ret.Get(0).(func(string) *core.Dogu); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Dogu)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *DoguRegistry) GetAll() ([]*core.Dogu, error) {
	ret := _m.Called()

	var r0 []*core.Dogu
	if rf, ok := ret.Get(0).(func() []*core.Dogu); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Dogu)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsEnabled provides a mock function with given fields: name
func (_m *DoguRegistry) IsEnabled(name string) (bool, error) {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: dogu
func (_m *DoguRegistry) Register(dogu *core.Dogu) error {
	ret := _m.Called(dogu)

	var r0 error
	if rf, ok := ret.Get(0).(func(*core.Dogu) error); ok {
		r0 = rf(dogu)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unregister provides a mock function with given fields: name
func (_m *DoguRegistry) Unregister(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

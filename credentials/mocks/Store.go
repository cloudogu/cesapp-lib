// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	"github.com/cloudogu/cesapp-lib/core"
	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Add provides a mock function with given fields: id, creds
func (_m *Store) Add(id string, creds *core.Credentials) error {
	ret := _m.Called(id, creds)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *core.Credentials) error); ok {
		r0 = rf(id, creds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *Store) Get(id string) *core.Credentials {
	ret := _m.Called(id)

	var r0 *core.Credentials
	if rf, ok := ret.Get(0).(func(string) *core.Credentials); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Credentials)
		}
	}

	return r0
}

// Remove provides a mock function with given fields: id
func (_m *Store) Remove(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockNewWatchConfigurationContextT is an autogenerated mock type for the NewWatchConfigurationContextT type
type MockNewWatchConfigurationContextT struct {
	mock.Mock
}

type MockNewWatchConfigurationContextT_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNewWatchConfigurationContextT) EXPECT() *MockNewWatchConfigurationContextT_Expecter {
	return &MockNewWatchConfigurationContextT_Expecter{mock: &_m.Mock}
}

// Cleanup provides a mock function with given fields: _a0
func (_m *MockNewWatchConfigurationContextT) Cleanup(_a0 func()) {
	_m.Called(_a0)
}

// MockNewWatchConfigurationContextT_Cleanup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Cleanup'
type MockNewWatchConfigurationContextT_Cleanup_Call struct {
	*mock.Call
}

// Cleanup is a helper method to define mock.On call
//  - _a0 func()
func (_e *MockNewWatchConfigurationContextT_Expecter) Cleanup(_a0 interface{}) *MockNewWatchConfigurationContextT_Cleanup_Call {
	return &MockNewWatchConfigurationContextT_Cleanup_Call{Call: _e.mock.On("Cleanup", _a0)}
}

func (_c *MockNewWatchConfigurationContextT_Cleanup_Call) Run(run func(_a0 func())) *MockNewWatchConfigurationContextT_Cleanup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func()))
	})
	return _c
}

func (_c *MockNewWatchConfigurationContextT_Cleanup_Call) Return() *MockNewWatchConfigurationContextT_Cleanup_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockNewWatchConfigurationContextT_Cleanup_Call) RunAndReturn(run func(func())) *MockNewWatchConfigurationContextT_Cleanup_Call {
	_c.Call.Return(run)
	return _c
}

// Errorf provides a mock function with given fields: format, args
func (_m *MockNewWatchConfigurationContextT) Errorf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockNewWatchConfigurationContextT_Errorf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Errorf'
type MockNewWatchConfigurationContextT_Errorf_Call struct {
	*mock.Call
}

// Errorf is a helper method to define mock.On call
//  - format string
//  - args ...interface{}
func (_e *MockNewWatchConfigurationContextT_Expecter) Errorf(format interface{}, args ...interface{}) *MockNewWatchConfigurationContextT_Errorf_Call {
	return &MockNewWatchConfigurationContextT_Errorf_Call{Call: _e.mock.On("Errorf",
		append([]interface{}{format}, args...)...)}
}

func (_c *MockNewWatchConfigurationContextT_Errorf_Call) Run(run func(format string, args ...interface{})) *MockNewWatchConfigurationContextT_Errorf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockNewWatchConfigurationContextT_Errorf_Call) Return() *MockNewWatchConfigurationContextT_Errorf_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockNewWatchConfigurationContextT_Errorf_Call) RunAndReturn(run func(string, ...interface{})) *MockNewWatchConfigurationContextT_Errorf_Call {
	_c.Call.Return(run)
	return _c
}

// FailNow provides a mock function with given fields:
func (_m *MockNewWatchConfigurationContextT) FailNow() {
	_m.Called()
}

// MockNewWatchConfigurationContextT_FailNow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FailNow'
type MockNewWatchConfigurationContextT_FailNow_Call struct {
	*mock.Call
}

// FailNow is a helper method to define mock.On call
func (_e *MockNewWatchConfigurationContextT_Expecter) FailNow() *MockNewWatchConfigurationContextT_FailNow_Call {
	return &MockNewWatchConfigurationContextT_FailNow_Call{Call: _e.mock.On("FailNow")}
}

func (_c *MockNewWatchConfigurationContextT_FailNow_Call) Run(run func()) *MockNewWatchConfigurationContextT_FailNow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNewWatchConfigurationContextT_FailNow_Call) Return() *MockNewWatchConfigurationContextT_FailNow_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockNewWatchConfigurationContextT_FailNow_Call) RunAndReturn(run func()) *MockNewWatchConfigurationContextT_FailNow_Call {
	_c.Call.Return(run)
	return _c
}

// Logf provides a mock function with given fields: format, args
func (_m *MockNewWatchConfigurationContextT) Logf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// MockNewWatchConfigurationContextT_Logf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Logf'
type MockNewWatchConfigurationContextT_Logf_Call struct {
	*mock.Call
}

// Logf is a helper method to define mock.On call
//  - format string
//  - args ...interface{}
func (_e *MockNewWatchConfigurationContextT_Expecter) Logf(format interface{}, args ...interface{}) *MockNewWatchConfigurationContextT_Logf_Call {
	return &MockNewWatchConfigurationContextT_Logf_Call{Call: _e.mock.On("Logf",
		append([]interface{}{format}, args...)...)}
}

func (_c *MockNewWatchConfigurationContextT_Logf_Call) Run(run func(format string, args ...interface{})) *MockNewWatchConfigurationContextT_Logf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockNewWatchConfigurationContextT_Logf_Call) Return() *MockNewWatchConfigurationContextT_Logf_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockNewWatchConfigurationContextT_Logf_Call) RunAndReturn(run func(string, ...interface{})) *MockNewWatchConfigurationContextT_Logf_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockNewWatchConfigurationContextT interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockNewWatchConfigurationContextT creates a new instance of MockNewWatchConfigurationContextT. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockNewWatchConfigurationContextT(t mockConstructorTestingTNewMockNewWatchConfigurationContextT) *MockNewWatchConfigurationContextT {
	mock := &MockNewWatchConfigurationContextT{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
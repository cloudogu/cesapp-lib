// Code generated by mockery v2.20.0. DO NOT EDIT.

package doguConf

import mock "github.com/stretchr/testify/mock"

// mockDoguConfigurationContext is an autogenerated mock type for the doguConfigurationContext type
type mockDoguConfigurationContext struct {
	mock.Mock
}

type mockDoguConfigurationContext_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguConfigurationContext) EXPECT() *mockDoguConfigurationContext_Expecter {
	return &mockDoguConfigurationContext_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: key
func (_m *mockDoguConfigurationContext) Delete(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigurationContext_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockDoguConfigurationContext_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//  - key string
func (_e *mockDoguConfigurationContext_Expecter) Delete(key interface{}) *mockDoguConfigurationContext_Delete_Call {
	return &mockDoguConfigurationContext_Delete_Call{Call: _e.mock.On("Delete", key)}
}

func (_c *mockDoguConfigurationContext_Delete_Call) Run(run func(key string)) *mockDoguConfigurationContext_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockDoguConfigurationContext_Delete_Call) Return(_a0 error) *mockDoguConfigurationContext_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigurationContext_Delete_Call) RunAndReturn(run func(string) error) *mockDoguConfigurationContext_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Exists provides a mock function with given fields: key
func (_m *mockDoguConfigurationContext) Exists(key string) (bool, error) {
	ret := _m.Called(key)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigurationContext_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type mockDoguConfigurationContext_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//  - key string
func (_e *mockDoguConfigurationContext_Expecter) Exists(key interface{}) *mockDoguConfigurationContext_Exists_Call {
	return &mockDoguConfigurationContext_Exists_Call{Call: _e.mock.On("Exists", key)}
}

func (_c *mockDoguConfigurationContext_Exists_Call) Run(run func(key string)) *mockDoguConfigurationContext_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockDoguConfigurationContext_Exists_Call) Return(_a0 bool, _a1 error) *mockDoguConfigurationContext_Exists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigurationContext_Exists_Call) RunAndReturn(run func(string) (bool, error)) *mockDoguConfigurationContext_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: key
func (_m *mockDoguConfigurationContext) Get(key string) (string, error) {
	ret := _m.Called(key)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigurationContext_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDoguConfigurationContext_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//  - key string
func (_e *mockDoguConfigurationContext_Expecter) Get(key interface{}) *mockDoguConfigurationContext_Get_Call {
	return &mockDoguConfigurationContext_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *mockDoguConfigurationContext_Get_Call) Run(run func(key string)) *mockDoguConfigurationContext_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockDoguConfigurationContext_Get_Call) Return(_a0 string, _a1 error) *mockDoguConfigurationContext_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigurationContext_Get_Call) RunAndReturn(run func(string) (string, error)) *mockDoguConfigurationContext_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: key, value
func (_m *mockDoguConfigurationContext) Set(key string, value string) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigurationContext_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type mockDoguConfigurationContext_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//  - key string
//  - value string
func (_e *mockDoguConfigurationContext_Expecter) Set(key interface{}, value interface{}) *mockDoguConfigurationContext_Set_Call {
	return &mockDoguConfigurationContext_Set_Call{Call: _e.mock.On("Set", key, value)}
}

func (_c *mockDoguConfigurationContext_Set_Call) Run(run func(key string, value string)) *mockDoguConfigurationContext_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *mockDoguConfigurationContext_Set_Call) Return(_a0 error) *mockDoguConfigurationContext_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigurationContext_Set_Call) RunAndReturn(run func(string, string) error) *mockDoguConfigurationContext_Set_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguConfigurationContext interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguConfigurationContext creates a new instance of mockDoguConfigurationContext. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguConfigurationContext(t mockConstructorTestingTnewMockDoguConfigurationContext) *mockDoguConfigurationContext {
	mock := &mockDoguConfigurationContext{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

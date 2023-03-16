// Code generated by mockery v2.20.0. DO NOT EDIT.

package registry

import mock "github.com/stretchr/testify/mock"

// MockState is an autogenerated mock type for the State type
type MockState struct {
	mock.Mock
}

type MockState_Expecter struct {
	mock *mock.Mock
}

func (_m *MockState) EXPECT() *MockState_Expecter {
	return &MockState_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields:
func (_m *MockState) Get() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockState_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockState_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
func (_e *MockState_Expecter) Get() *MockState_Get_Call {
	return &MockState_Get_Call{Call: _e.mock.On("Get")}
}

func (_c *MockState_Get_Call) Run(run func()) *MockState_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockState_Get_Call) Return(_a0 string, _a1 error) *MockState_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockState_Get_Call) RunAndReturn(run func() (string, error)) *MockState_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields:
func (_m *MockState) Remove() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockState_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type MockState_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
func (_e *MockState_Expecter) Remove() *MockState_Remove_Call {
	return &MockState_Remove_Call{Call: _e.mock.On("Remove")}
}

func (_c *MockState_Remove_Call) Run(run func()) *MockState_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockState_Remove_Call) Return(_a0 error) *MockState_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockState_Remove_Call) RunAndReturn(run func() error) *MockState_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: value
func (_m *MockState) Set(value string) error {
	ret := _m.Called(value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockState_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockState_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//  - value string
func (_e *MockState_Expecter) Set(value interface{}) *MockState_Set_Call {
	return &MockState_Set_Call{Call: _e.mock.On("Set", value)}
}

func (_c *MockState_Set_Call) Run(run func(value string)) *MockState_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockState_Set_Call) Return(_a0 error) *MockState_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockState_Set_Call) RunAndReturn(run func(string) error) *MockState_Set_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockState interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockState creates a new instance of MockState. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockState(t mockConstructorTestingTNewMockState) *MockState {
	mock := &MockState{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

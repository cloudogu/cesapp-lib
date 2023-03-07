// Code generated by mockery v2.20.0. DO NOT EDIT.

package remote

import (
	core "github.com/cloudogu/cesapp-lib/core"
	mock "github.com/stretchr/testify/mock"
)

// MockRegistry is an autogenerated mock type for the Registry type
type MockRegistry struct {
	mock.Mock
}

type MockRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRegistry) EXPECT() *MockRegistry_Expecter {
	return &MockRegistry_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: dogu
func (_m *MockRegistry) Create(dogu *core.Dogu) error {
	ret := _m.Called(dogu)

	var r0 error
	if rf, ok := ret.Get(0).(func(*core.Dogu) error); ok {
		r0 = rf(dogu)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRegistry_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockRegistry_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//  - dogu *core.Dogu
func (_e *MockRegistry_Expecter) Create(dogu interface{}) *MockRegistry_Create_Call {
	return &MockRegistry_Create_Call{Call: _e.mock.On("Create", dogu)}
}

func (_c *MockRegistry_Create_Call) Run(run func(dogu *core.Dogu)) *MockRegistry_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*core.Dogu))
	})
	return _c
}

func (_c *MockRegistry_Create_Call) Return(_a0 error) *MockRegistry_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRegistry_Create_Call) RunAndReturn(run func(*core.Dogu) error) *MockRegistry_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: dogu
func (_m *MockRegistry) Delete(dogu *core.Dogu) error {
	ret := _m.Called(dogu)

	var r0 error
	if rf, ok := ret.Get(0).(func(*core.Dogu) error); ok {
		r0 = rf(dogu)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRegistry_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockRegistry_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//  - dogu *core.Dogu
func (_e *MockRegistry_Expecter) Delete(dogu interface{}) *MockRegistry_Delete_Call {
	return &MockRegistry_Delete_Call{Call: _e.mock.On("Delete", dogu)}
}

func (_c *MockRegistry_Delete_Call) Run(run func(dogu *core.Dogu)) *MockRegistry_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*core.Dogu))
	})
	return _c
}

func (_c *MockRegistry_Delete_Call) Return(_a0 error) *MockRegistry_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRegistry_Delete_Call) RunAndReturn(run func(*core.Dogu) error) *MockRegistry_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: name
func (_m *MockRegistry) Get(name string) (*core.Dogu, error) {
	ret := _m.Called(name)

	var r0 *core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*core.Dogu, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *core.Dogu); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRegistry_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockRegistry_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//  - name string
func (_e *MockRegistry_Expecter) Get(name interface{}) *MockRegistry_Get_Call {
	return &MockRegistry_Get_Call{Call: _e.mock.On("Get", name)}
}

func (_c *MockRegistry_Get_Call) Run(run func(name string)) *MockRegistry_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRegistry_Get_Call) Return(_a0 *core.Dogu, _a1 error) *MockRegistry_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRegistry_Get_Call) RunAndReturn(run func(string) (*core.Dogu, error)) *MockRegistry_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields:
func (_m *MockRegistry) GetAll() ([]*core.Dogu, error) {
	ret := _m.Called()

	var r0 []*core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*core.Dogu, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*core.Dogu); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRegistry_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockRegistry_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *MockRegistry_Expecter) GetAll() *MockRegistry_GetAll_Call {
	return &MockRegistry_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *MockRegistry_GetAll_Call) Run(run func()) *MockRegistry_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRegistry_GetAll_Call) Return(_a0 []*core.Dogu, _a1 error) *MockRegistry_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRegistry_GetAll_Call) RunAndReturn(run func() ([]*core.Dogu, error)) *MockRegistry_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetVersion provides a mock function with given fields: name, version
func (_m *MockRegistry) GetVersion(name string, version string) (*core.Dogu, error) {
	ret := _m.Called(name, version)

	var r0 *core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*core.Dogu, error)); ok {
		return rf(name, version)
	}
	if rf, ok := ret.Get(0).(func(string, string) *core.Dogu); ok {
		r0 = rf(name, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(name, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRegistry_GetVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVersion'
type MockRegistry_GetVersion_Call struct {
	*mock.Call
}

// GetVersion is a helper method to define mock.On call
//  - name string
//  - version string
func (_e *MockRegistry_Expecter) GetVersion(name interface{}, version interface{}) *MockRegistry_GetVersion_Call {
	return &MockRegistry_GetVersion_Call{Call: _e.mock.On("GetVersion", name, version)}
}

func (_c *MockRegistry_GetVersion_Call) Run(run func(name string, version string)) *MockRegistry_GetVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockRegistry_GetVersion_Call) Return(_a0 *core.Dogu, _a1 error) *MockRegistry_GetVersion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRegistry_GetVersion_Call) RunAndReturn(run func(string, string) (*core.Dogu, error)) *MockRegistry_GetVersion_Call {
	_c.Call.Return(run)
	return _c
}

// GetVersionsOf provides a mock function with given fields: name
func (_m *MockRegistry) GetVersionsOf(name string) ([]core.Version, error) {
	ret := _m.Called(name)

	var r0 []core.Version
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]core.Version, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) []core.Version); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Version)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRegistry_GetVersionsOf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVersionsOf'
type MockRegistry_GetVersionsOf_Call struct {
	*mock.Call
}

// GetVersionsOf is a helper method to define mock.On call
//  - name string
func (_e *MockRegistry_Expecter) GetVersionsOf(name interface{}) *MockRegistry_GetVersionsOf_Call {
	return &MockRegistry_GetVersionsOf_Call{Call: _e.mock.On("GetVersionsOf", name)}
}

func (_c *MockRegistry_GetVersionsOf_Call) Run(run func(name string)) *MockRegistry_GetVersionsOf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRegistry_GetVersionsOf_Call) Return(_a0 []core.Version, _a1 error) *MockRegistry_GetVersionsOf_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRegistry_GetVersionsOf_Call) RunAndReturn(run func(string) ([]core.Version, error)) *MockRegistry_GetVersionsOf_Call {
	_c.Call.Return(run)
	return _c
}

// SetUseCache provides a mock function with given fields: useCache
func (_m *MockRegistry) SetUseCache(useCache bool) {
	_m.Called(useCache)
}

// MockRegistry_SetUseCache_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetUseCache'
type MockRegistry_SetUseCache_Call struct {
	*mock.Call
}

// SetUseCache is a helper method to define mock.On call
//  - useCache bool
func (_e *MockRegistry_Expecter) SetUseCache(useCache interface{}) *MockRegistry_SetUseCache_Call {
	return &MockRegistry_SetUseCache_Call{Call: _e.mock.On("SetUseCache", useCache)}
}

func (_c *MockRegistry_SetUseCache_Call) Run(run func(useCache bool)) *MockRegistry_SetUseCache_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *MockRegistry_SetUseCache_Call) Return() *MockRegistry_SetUseCache_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockRegistry_SetUseCache_Call) RunAndReturn(run func(bool)) *MockRegistry_SetUseCache_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockRegistry interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRegistry creates a new instance of MockRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRegistry(t mockConstructorTestingTNewMockRegistry) *MockRegistry {
	mock := &MockRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
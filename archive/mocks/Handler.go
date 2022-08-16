// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Handler is an autogenerated mock type for the Handler type
type Handler struct {
	mock.Mock
}

// AppendFileToArchive provides a mock function with given fields: fileToZipPath, filepathInZip
func (_m *Handler) AppendFileToArchive(fileToZipPath string, filepathInZip string) error {
	ret := _m.Called(fileToZipPath, filepathInZip)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(fileToZipPath, filepathInZip)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Handler) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateZipArchiveFile provides a mock function with given fields: zipFilePath
func (_m *Handler) CreateZipArchiveFile(zipFilePath string) (io.Writer, error) {
	ret := _m.Called(zipFilePath)

	var r0 io.Writer
	if rf, ok := ret.Get(0).(func(string) io.Writer); ok {
		r0 = rf(zipFilePath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Writer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(zipFilePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitialiseZipWriter provides a mock function with given fields: zipFile
func (_m *Handler) InitializeZipWriter(zipFile io.Writer) {
	_m.Called(zipFile)
}

// WriteFilesIntoArchive provides a mock function with given fields: filePaths, closeAfterFinish
func (_m *Handler) WriteFilesIntoArchive(filePaths []string, closeAfterFinish bool) error {
	ret := _m.Called(filePaths, closeAfterFinish)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, bool) error); ok {
		r0 = rf(filePaths, closeAfterFinish)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewHandler creates a new instance of Handler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHandler(t mockConstructorTestingTNewHandler) *Handler {
	mock := &Handler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

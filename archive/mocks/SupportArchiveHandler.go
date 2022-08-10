// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// SupportArchiveHandler is an autogenerated mock type for the SupportArchiveHandler type
type SupportArchiveHandler struct {
	mock.Mock
}

// AppendFileToArchive provides a mock function with given fields: fileToZipPath, filepathInZip
func (_m *SupportArchiveHandler) AppendFileToArchive(fileToZipPath string, filepathInZip string) {
	_m.Called(fileToZipPath, filepathInZip)
}

// Close provides a mock function with given fields:
func (_m *SupportArchiveHandler) Close() error {
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
func (_m *SupportArchiveHandler) CreateZipArchiveFile(zipFilePath string) (io.Writer, error) {
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
func (_m *SupportArchiveHandler) InitialiseZipWriter(zipFile io.Writer) {
	_m.Called(zipFile)
}

// WriteFilesIntoArchive provides a mock function with given fields: filePaths, closeAfterFinish
func (_m *SupportArchiveHandler) WriteFilesIntoArchive(filePaths []string, closeAfterFinish bool) {
	_m.Called(filePaths, closeAfterFinish)
}

type mockConstructorTestingTNewSupportArchiveHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewSupportArchiveHandler creates a new instance of SupportArchiveHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSupportArchiveHandler(t mockConstructorTestingTNewSupportArchiveHandler) *SupportArchiveHandler {
	mock := &SupportArchiveHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

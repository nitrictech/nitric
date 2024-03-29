// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azure/azure-storage-blob-go/azblob (interfaces: StorageError)

// Package mock_azblob is a generated GoMock package.
package mock_azblob

import (
	http "net/http"
	reflect "reflect"

	azblob "github.com/Azure/azure-storage-blob-go/azblob"
	gomock "github.com/golang/mock/gomock"
)

// MockStorageError is a mock of StorageError interface.
type MockStorageError struct {
	ctrl     *gomock.Controller
	recorder *MockStorageErrorMockRecorder
}

// MockStorageErrorMockRecorder is the mock recorder for MockStorageError.
type MockStorageErrorMockRecorder struct {
	mock *MockStorageError
}

// NewMockStorageError creates a new mock instance.
func NewMockStorageError(ctrl *gomock.Controller) *MockStorageError {
	mock := &MockStorageError{ctrl: ctrl}
	mock.recorder = &MockStorageErrorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageError) EXPECT() *MockStorageErrorMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockStorageError) Error() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(string)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockStorageErrorMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockStorageError)(nil).Error))
}

// Response mocks base method.
func (m *MockStorageError) Response() *http.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Response")
	ret0, _ := ret[0].(*http.Response)
	return ret0
}

// Response indicates an expected call of Response.
func (mr *MockStorageErrorMockRecorder) Response() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Response", reflect.TypeOf((*MockStorageError)(nil).Response))
}

// ServiceCode mocks base method.
func (m *MockStorageError) ServiceCode() azblob.ServiceCodeType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceCode")
	ret0, _ := ret[0].(azblob.ServiceCodeType)
	return ret0
}

// ServiceCode indicates an expected call of ServiceCode.
func (mr *MockStorageErrorMockRecorder) ServiceCode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceCode", reflect.TypeOf((*MockStorageError)(nil).ServiceCode))
}

// Temporary mocks base method.
func (m *MockStorageError) Temporary() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Temporary")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Temporary indicates an expected call of Temporary.
func (mr *MockStorageErrorMockRecorder) Temporary() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Temporary", reflect.TypeOf((*MockStorageError)(nil).Temporary))
}

// Timeout mocks base method.
func (m *MockStorageError) Timeout() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timeout")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Timeout indicates an expected call of Timeout.
func (mr *MockStorageErrorMockRecorder) Timeout() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timeout", reflect.TypeOf((*MockStorageError)(nil).Timeout))
}

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/nitric/core/pkg/workers/http (interfaces: HttpRequestHandler)

// Package mock_http is a generated GoMock package.
package mock_http

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	fasthttp "github.com/valyala/fasthttp"
)

// MockHttpRequestHandler is a mock of HttpRequestHandler interface.
type MockHttpRequestHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHttpRequestHandlerMockRecorder
}

// MockHttpRequestHandlerMockRecorder is the mock recorder for MockHttpRequestHandler.
type MockHttpRequestHandlerMockRecorder struct {
	mock *MockHttpRequestHandler
}

// NewMockHttpRequestHandler creates a new mock instance.
func NewMockHttpRequestHandler(ctrl *gomock.Controller) *MockHttpRequestHandler {
	mock := &MockHttpRequestHandler{ctrl: ctrl}
	mock.recorder = &MockHttpRequestHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpRequestHandler) EXPECT() *MockHttpRequestHandlerMockRecorder {
	return m.recorder
}

// HandleRequest mocks base method.
func (m *MockHttpRequestHandler) HandleRequest(arg0 *fasthttp.Request) (*fasthttp.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleRequest", arg0)
	ret0, _ := ret[0].(*fasthttp.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleRequest indicates an expected call of HandleRequest.
func (mr *MockHttpRequestHandlerMockRecorder) HandleRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleRequest", reflect.TypeOf((*MockHttpRequestHandler)(nil).HandleRequest), arg0)
}

// WorkerCount mocks base method.
func (m *MockHttpRequestHandler) WorkerCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WorkerCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// WorkerCount indicates an expected call of WorkerCount.
func (mr *MockHttpRequestHandlerMockRecorder) WorkerCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WorkerCount", reflect.TypeOf((*MockHttpRequestHandler)(nil).WorkerCount))
}
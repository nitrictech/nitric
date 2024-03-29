// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/nitric/core/pkg/workers/schedules (interfaces: ScheduleRequestHandler)

// Package mock_schedules is a generated GoMock package.
package mock_schedules

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
)

// MockScheduleRequestHandler is a mock of ScheduleRequestHandler interface.
type MockScheduleRequestHandler struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleRequestHandlerMockRecorder
}

// MockScheduleRequestHandlerMockRecorder is the mock recorder for MockScheduleRequestHandler.
type MockScheduleRequestHandlerMockRecorder struct {
	mock *MockScheduleRequestHandler
}

// NewMockScheduleRequestHandler creates a new mock instance.
func NewMockScheduleRequestHandler(ctrl *gomock.Controller) *MockScheduleRequestHandler {
	mock := &MockScheduleRequestHandler{ctrl: ctrl}
	mock.recorder = &MockScheduleRequestHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduleRequestHandler) EXPECT() *MockScheduleRequestHandlerMockRecorder {
	return m.recorder
}

// HandleRequest mocks base method.
func (m *MockScheduleRequestHandler) HandleRequest(arg0 *schedulespb.ServerMessage) (*schedulespb.ClientMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleRequest", arg0)
	ret0, _ := ret[0].(*schedulespb.ClientMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleRequest indicates an expected call of HandleRequest.
func (mr *MockScheduleRequestHandlerMockRecorder) HandleRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleRequest", reflect.TypeOf((*MockScheduleRequestHandler)(nil).HandleRequest), arg0)
}

// Schedule mocks base method.
func (m *MockScheduleRequestHandler) Schedule(arg0 schedulespb.Schedules_ScheduleServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Schedule", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Schedule indicates an expected call of Schedule.
func (mr *MockScheduleRequestHandlerMockRecorder) Schedule(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Schedule", reflect.TypeOf((*MockScheduleRequestHandler)(nil).Schedule), arg0)
}

// WorkerCount mocks base method.
func (m *MockScheduleRequestHandler) WorkerCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WorkerCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// WorkerCount indicates an expected call of WorkerCount.
func (mr *MockScheduleRequestHandlerMockRecorder) WorkerCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WorkerCount", reflect.TypeOf((*MockScheduleRequestHandler)(nil).WorkerCount))
}

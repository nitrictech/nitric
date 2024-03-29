// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub (interfaces: PubsubClient,TopicIterator,Topic,PublishResult)

// Package mock_pubsub is a generated GoMock package.
package mock_pubsub

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ifaces_pubsub "github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub"
)

// MockPubsubClient is a mock of PubsubClient interface.
type MockPubsubClient struct {
	ctrl     *gomock.Controller
	recorder *MockPubsubClientMockRecorder
}

// MockPubsubClientMockRecorder is the mock recorder for MockPubsubClient.
type MockPubsubClientMockRecorder struct {
	mock *MockPubsubClient
}

// NewMockPubsubClient creates a new mock instance.
func NewMockPubsubClient(ctrl *gomock.Controller) *MockPubsubClient {
	mock := &MockPubsubClient{ctrl: ctrl}
	mock.recorder = &MockPubsubClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubsubClient) EXPECT() *MockPubsubClientMockRecorder {
	return m.recorder
}

// Topic mocks base method.
func (m *MockPubsubClient) Topic(arg0 string) ifaces_pubsub.Topic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Topic", arg0)
	ret0, _ := ret[0].(ifaces_pubsub.Topic)
	return ret0
}

// Topic indicates an expected call of Topic.
func (mr *MockPubsubClientMockRecorder) Topic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Topic", reflect.TypeOf((*MockPubsubClient)(nil).Topic), arg0)
}

// Topics mocks base method.
func (m *MockPubsubClient) Topics(arg0 context.Context) ifaces_pubsub.TopicIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Topics", arg0)
	ret0, _ := ret[0].(ifaces_pubsub.TopicIterator)
	return ret0
}

// Topics indicates an expected call of Topics.
func (mr *MockPubsubClientMockRecorder) Topics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Topics", reflect.TypeOf((*MockPubsubClient)(nil).Topics), arg0)
}

// MockTopicIterator is a mock of TopicIterator interface.
type MockTopicIterator struct {
	ctrl     *gomock.Controller
	recorder *MockTopicIteratorMockRecorder
}

// MockTopicIteratorMockRecorder is the mock recorder for MockTopicIterator.
type MockTopicIteratorMockRecorder struct {
	mock *MockTopicIterator
}

// NewMockTopicIterator creates a new mock instance.
func NewMockTopicIterator(ctrl *gomock.Controller) *MockTopicIterator {
	mock := &MockTopicIterator{ctrl: ctrl}
	mock.recorder = &MockTopicIteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopicIterator) EXPECT() *MockTopicIteratorMockRecorder {
	return m.recorder
}

// Next mocks base method.
func (m *MockTopicIterator) Next() (ifaces_pubsub.Topic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(ifaces_pubsub.Topic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Next indicates an expected call of Next.
func (mr *MockTopicIteratorMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockTopicIterator)(nil).Next))
}

// MockTopic is a mock of Topic interface.
type MockTopic struct {
	ctrl     *gomock.Controller
	recorder *MockTopicMockRecorder
}

// MockTopicMockRecorder is the mock recorder for MockTopic.
type MockTopicMockRecorder struct {
	mock *MockTopic
}

// NewMockTopic creates a new mock instance.
func NewMockTopic(ctrl *gomock.Controller) *MockTopic {
	mock := &MockTopic{ctrl: ctrl}
	mock.recorder = &MockTopicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopic) EXPECT() *MockTopicMockRecorder {
	return m.recorder
}

// Exists mocks base method.
func (m *MockTopic) Exists(arg0 context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockTopicMockRecorder) Exists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockTopic)(nil).Exists), arg0)
}

// ID mocks base method.
func (m *MockTopic) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockTopicMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockTopic)(nil).ID))
}

// Labels mocks base method.
func (m *MockTopic) Labels(arg0 context.Context) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Labels", arg0)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Labels indicates an expected call of Labels.
func (mr *MockTopicMockRecorder) Labels(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Labels", reflect.TypeOf((*MockTopic)(nil).Labels), arg0)
}

// Publish mocks base method.
func (m *MockTopic) Publish(arg0 context.Context, arg1 ifaces_pubsub.Message) ifaces_pubsub.PublishResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(ifaces_pubsub.PublishResult)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockTopicMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockTopic)(nil).Publish), arg0, arg1)
}

// String mocks base method.
func (m *MockTopic) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockTopicMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockTopic)(nil).String))
}

// Subscriptions mocks base method.
func (m *MockTopic) Subscriptions(arg0 context.Context) ifaces_pubsub.SubscriptionIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscriptions", arg0)
	ret0, _ := ret[0].(ifaces_pubsub.SubscriptionIterator)
	return ret0
}

// Subscriptions indicates an expected call of Subscriptions.
func (mr *MockTopicMockRecorder) Subscriptions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscriptions", reflect.TypeOf((*MockTopic)(nil).Subscriptions), arg0)
}

// MockPublishResult is a mock of PublishResult interface.
type MockPublishResult struct {
	ctrl     *gomock.Controller
	recorder *MockPublishResultMockRecorder
}

// MockPublishResultMockRecorder is the mock recorder for MockPublishResult.
type MockPublishResultMockRecorder struct {
	mock *MockPublishResult
}

// NewMockPublishResult creates a new mock instance.
func NewMockPublishResult(ctrl *gomock.Controller) *MockPublishResult {
	mock := &MockPublishResult{ctrl: ctrl}
	mock.recorder = &MockPublishResultMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPublishResult) EXPECT() *MockPublishResultMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPublishResult) Get(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPublishResultMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPublishResult)(nil).Get), arg0)
}

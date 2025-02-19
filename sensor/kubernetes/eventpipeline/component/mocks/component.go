// Code generated by MockGen. DO NOT EDIT.
// Source: component.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	message "github.com/stackrox/rox/sensor/common/message"
	component "github.com/stackrox/rox/sensor/kubernetes/eventpipeline/component"
	gomock "go.uber.org/mock/gomock"
)

// MockPipelineComponent is a mock of PipelineComponent interface.
type MockPipelineComponent struct {
	ctrl     *gomock.Controller
	recorder *MockPipelineComponentMockRecorder
}

// MockPipelineComponentMockRecorder is the mock recorder for MockPipelineComponent.
type MockPipelineComponentMockRecorder struct {
	mock *MockPipelineComponent
}

// NewMockPipelineComponent creates a new mock instance.
func NewMockPipelineComponent(ctrl *gomock.Controller) *MockPipelineComponent {
	mock := &MockPipelineComponent{ctrl: ctrl}
	mock.recorder = &MockPipelineComponentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPipelineComponent) EXPECT() *MockPipelineComponentMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockPipelineComponent) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockPipelineComponentMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockPipelineComponent)(nil).Start))
}

// Stop mocks base method.
func (m *MockPipelineComponent) Stop(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop", arg0)
}

// Stop indicates an expected call of Stop.
func (mr *MockPipelineComponentMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockPipelineComponent)(nil).Stop), arg0)
}

// MockResolver is a mock of Resolver interface.
type MockResolver struct {
	ctrl     *gomock.Controller
	recorder *MockResolverMockRecorder
}

// MockResolverMockRecorder is the mock recorder for MockResolver.
type MockResolverMockRecorder struct {
	mock *MockResolver
}

// NewMockResolver creates a new mock instance.
func NewMockResolver(ctrl *gomock.Controller) *MockResolver {
	mock := &MockResolver{ctrl: ctrl}
	mock.recorder = &MockResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResolver) EXPECT() *MockResolverMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockResolver) Send(event *component.ResourceEvent) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Send", event)
}

// Send indicates an expected call of Send.
func (mr *MockResolverMockRecorder) Send(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockResolver)(nil).Send), event)
}

// Start mocks base method.
func (m *MockResolver) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockResolverMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockResolver)(nil).Start))
}

// Stop mocks base method.
func (m *MockResolver) Stop(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop", arg0)
}

// Stop indicates an expected call of Stop.
func (mr *MockResolverMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockResolver)(nil).Stop), arg0)
}

// MockOutputQueue is a mock of OutputQueue interface.
type MockOutputQueue struct {
	ctrl     *gomock.Controller
	recorder *MockOutputQueueMockRecorder
}

// MockOutputQueueMockRecorder is the mock recorder for MockOutputQueue.
type MockOutputQueueMockRecorder struct {
	mock *MockOutputQueue
}

// NewMockOutputQueue creates a new mock instance.
func NewMockOutputQueue(ctrl *gomock.Controller) *MockOutputQueue {
	mock := &MockOutputQueue{ctrl: ctrl}
	mock.recorder = &MockOutputQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOutputQueue) EXPECT() *MockOutputQueueMockRecorder {
	return m.recorder
}

// ResponsesC mocks base method.
func (m *MockOutputQueue) ResponsesC() <-chan *message.ExpiringMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResponsesC")
	ret0, _ := ret[0].(<-chan *message.ExpiringMessage)
	return ret0
}

// ResponsesC indicates an expected call of ResponsesC.
func (mr *MockOutputQueueMockRecorder) ResponsesC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponsesC", reflect.TypeOf((*MockOutputQueue)(nil).ResponsesC))
}

// Send mocks base method.
func (m *MockOutputQueue) Send(event *component.ResourceEvent) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Send", event)
}

// Send indicates an expected call of Send.
func (mr *MockOutputQueueMockRecorder) Send(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockOutputQueue)(nil).Send), event)
}

// Start mocks base method.
func (m *MockOutputQueue) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockOutputQueueMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockOutputQueue)(nil).Start))
}

// Stop mocks base method.
func (m *MockOutputQueue) Stop(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop", arg0)
}

// Stop indicates an expected call of Stop.
func (mr *MockOutputQueueMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockOutputQueue)(nil).Stop), arg0)
}

// MockContextListener is a mock of ContextListener interface.
type MockContextListener struct {
	ctrl     *gomock.Controller
	recorder *MockContextListenerMockRecorder
}

// MockContextListenerMockRecorder is the mock recorder for MockContextListener.
type MockContextListenerMockRecorder struct {
	mock *MockContextListener
}

// NewMockContextListener creates a new mock instance.
func NewMockContextListener(ctrl *gomock.Controller) *MockContextListener {
	mock := &MockContextListener{ctrl: ctrl}
	mock.recorder = &MockContextListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContextListener) EXPECT() *MockContextListenerMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockContextListener) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockContextListenerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockContextListener)(nil).Start))
}

// StartWithContext mocks base method.
func (m *MockContextListener) StartWithContext(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartWithContext", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartWithContext indicates an expected call of StartWithContext.
func (mr *MockContextListenerMockRecorder) StartWithContext(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartWithContext", reflect.TypeOf((*MockContextListener)(nil).StartWithContext), arg0)
}

// Stop mocks base method.
func (m *MockContextListener) Stop(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop", arg0)
}

// Stop indicates an expected call of Stop.
func (mr *MockContextListenerMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockContextListener)(nil).Stop), arg0)
}

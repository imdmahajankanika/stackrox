// Code generated by MockGen. DO NOT EDIT.
// Source: flow.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	types "github.com/gogo/protobuf/types"
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	timestamp "github.com/stackrox/rox/pkg/timestamp"
)

// MockFlowStore is a mock of FlowStore interface.
type MockFlowStore struct {
	ctrl     *gomock.Controller
	recorder *MockFlowStoreMockRecorder
}

// MockFlowStoreMockRecorder is the mock recorder for MockFlowStore.
type MockFlowStoreMockRecorder struct {
	mock *MockFlowStore
}

// NewMockFlowStore creates a new mock instance.
func NewMockFlowStore(ctrl *gomock.Controller) *MockFlowStore {
	mock := &MockFlowStore{ctrl: ctrl}
	mock.recorder = &MockFlowStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlowStore) EXPECT() *MockFlowStoreMockRecorder {
	return m.recorder
}

// GetAllFlows mocks base method.
func (m *MockFlowStore) GetAllFlows(ctx context.Context, since *types.Timestamp) ([]*storage.NetworkFlow, *types.Timestamp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllFlows", ctx, since)
	ret0, _ := ret[0].([]*storage.NetworkFlow)
	ret1, _ := ret[1].(*types.Timestamp)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllFlows indicates an expected call of GetAllFlows.
func (mr *MockFlowStoreMockRecorder) GetAllFlows(ctx, since interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllFlows", reflect.TypeOf((*MockFlowStore)(nil).GetAllFlows), ctx, since)
}

// GetFlowsForDeployment mocks base method.
func (m *MockFlowStore) GetFlowsForDeployment(ctx context.Context, deploymentID string) ([]*storage.NetworkFlow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowsForDeployment", ctx, deploymentID)
	ret0, _ := ret[0].([]*storage.NetworkFlow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlowsForDeployment indicates an expected call of GetFlowsForDeployment.
func (mr *MockFlowStoreMockRecorder) GetFlowsForDeployment(ctx, deploymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowsForDeployment", reflect.TypeOf((*MockFlowStore)(nil).GetFlowsForDeployment), ctx, deploymentID)
}

// GetMatchingFlows mocks base method.
func (m *MockFlowStore) GetMatchingFlows(ctx context.Context, pred func(*storage.NetworkFlowProperties) bool, since *types.Timestamp) ([]*storage.NetworkFlow, *types.Timestamp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchingFlows", ctx, pred, since)
	ret0, _ := ret[0].([]*storage.NetworkFlow)
	ret1, _ := ret[1].(*types.Timestamp)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMatchingFlows indicates an expected call of GetMatchingFlows.
func (mr *MockFlowStoreMockRecorder) GetMatchingFlows(ctx, pred, since interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchingFlows", reflect.TypeOf((*MockFlowStore)(nil).GetMatchingFlows), ctx, pred, since)
}

// RemoveFlow mocks base method.
func (m *MockFlowStore) RemoveFlow(ctx context.Context, props *storage.NetworkFlowProperties) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFlow", ctx, props)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFlow indicates an expected call of RemoveFlow.
func (mr *MockFlowStoreMockRecorder) RemoveFlow(ctx, props interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFlow", reflect.TypeOf((*MockFlowStore)(nil).RemoveFlow), ctx, props)
}

// RemoveFlowsForDeployment mocks base method.
func (m *MockFlowStore) RemoveFlowsForDeployment(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFlowsForDeployment", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFlowsForDeployment indicates an expected call of RemoveFlowsForDeployment.
func (mr *MockFlowStoreMockRecorder) RemoveFlowsForDeployment(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFlowsForDeployment", reflect.TypeOf((*MockFlowStore)(nil).RemoveFlowsForDeployment), ctx, id)
}

// RemoveOrphanedFlows mocks base method.
func (m *MockFlowStore) RemoveOrphanedFlows(ctx context.Context, orphanWindow *time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOrphanedFlows", ctx, orphanWindow)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveOrphanedFlows indicates an expected call of RemoveOrphanedFlows.
func (mr *MockFlowStoreMockRecorder) RemoveOrphanedFlows(ctx, orphanWindow interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOrphanedFlows", reflect.TypeOf((*MockFlowStore)(nil).RemoveOrphanedFlows), ctx, orphanWindow)
}

// RemoveStaleFlows mocks base method.
func (m *MockFlowStore) RemoveStaleFlows(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveStaleFlows", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveStaleFlows indicates an expected call of RemoveStaleFlows.
func (mr *MockFlowStoreMockRecorder) RemoveStaleFlows(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveStaleFlows", reflect.TypeOf((*MockFlowStore)(nil).RemoveStaleFlows), ctx)
}

// UpsertFlows mocks base method.
func (m *MockFlowStore) UpsertFlows(ctx context.Context, flows []*storage.NetworkFlow, lastUpdateTS timestamp.MicroTS) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertFlows", ctx, flows, lastUpdateTS)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertFlows indicates an expected call of UpsertFlows.
func (mr *MockFlowStoreMockRecorder) UpsertFlows(ctx, flows, lastUpdateTS interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertFlows", reflect.TypeOf((*MockFlowStore)(nil).UpsertFlows), ctx, flows, lastUpdateTS)
}

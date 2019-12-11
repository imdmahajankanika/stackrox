// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/compliance/datastore (interfaces: SacFilter)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	compliance "github.com/stackrox/rox/central/compliance"
	types "github.com/stackrox/rox/central/compliance/datastore/types"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
)

// MockSacFilter is a mock of SacFilter interface
type MockSacFilter struct {
	ctrl     *gomock.Controller
	recorder *MockSacFilterMockRecorder
}

// MockSacFilterMockRecorder is the mock recorder for MockSacFilter
type MockSacFilterMockRecorder struct {
	mock *MockSacFilter
}

// NewMockSacFilter creates a new mock instance
func NewMockSacFilter(ctrl *gomock.Controller) *MockSacFilter {
	mock := &MockSacFilter{ctrl: ctrl}
	mock.recorder = &MockSacFilterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSacFilter) EXPECT() *MockSacFilterMockRecorder {
	return m.recorder
}

// FilterBatchResults mocks base method
func (m *MockSacFilter) FilterBatchResults(arg0 context.Context, arg1 map[compliance.ClusterStandardPair]types.ResultsWithStatus) (map[compliance.ClusterStandardPair]types.ResultsWithStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterBatchResults", arg0, arg1)
	ret0, _ := ret[0].(map[compliance.ClusterStandardPair]types.ResultsWithStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterBatchResults indicates an expected call of FilterBatchResults
func (mr *MockSacFilterMockRecorder) FilterBatchResults(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterBatchResults", reflect.TypeOf((*MockSacFilter)(nil).FilterBatchResults), arg0, arg1)
}

// FilterRunResults mocks base method
func (m *MockSacFilter) FilterRunResults(arg0 context.Context, arg1 *storage.ComplianceRunResults) (*storage.ComplianceRunResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterRunResults", arg0, arg1)
	ret0, _ := ret[0].(*storage.ComplianceRunResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterRunResults indicates an expected call of FilterRunResults
func (mr *MockSacFilterMockRecorder) FilterRunResults(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterRunResults", reflect.TypeOf((*MockSacFilter)(nil).FilterRunResults), arg0, arg1)
}

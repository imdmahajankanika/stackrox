// Code generated by MockGen. DO NOT EDIT.
// Source: indexer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

// MockIndexer is a mock of Indexer interface.
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer.
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance.
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// AddListAlert mocks base method.
func (m *MockIndexer) AddListAlert(listalert *storage.ListAlert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddListAlert", listalert)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddListAlert indicates an expected call of AddListAlert.
func (mr *MockIndexerMockRecorder) AddListAlert(listalert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddListAlert", reflect.TypeOf((*MockIndexer)(nil).AddListAlert), listalert)
}

// AddListAlerts mocks base method.
func (m *MockIndexer) AddListAlerts(listalerts []*storage.ListAlert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddListAlerts", listalerts)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddListAlerts indicates an expected call of AddListAlerts.
func (mr *MockIndexerMockRecorder) AddListAlerts(listalerts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddListAlerts", reflect.TypeOf((*MockIndexer)(nil).AddListAlerts), listalerts)
}

// Count mocks base method.
func (m *MockIndexer) Count(ctx context.Context, q *v1.Query, opts ...blevesearch.SearchOption) (int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, q}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Count", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockIndexerMockRecorder) Count(ctx, q interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, q}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockIndexer)(nil).Count), varargs...)
}

// DeleteListAlert mocks base method.
func (m *MockIndexer) DeleteListAlert(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteListAlert", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteListAlert indicates an expected call of DeleteListAlert.
func (mr *MockIndexerMockRecorder) DeleteListAlert(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteListAlert", reflect.TypeOf((*MockIndexer)(nil).DeleteListAlert), id)
}

// DeleteListAlerts mocks base method.
func (m *MockIndexer) DeleteListAlerts(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteListAlerts", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteListAlerts indicates an expected call of DeleteListAlerts.
func (mr *MockIndexerMockRecorder) DeleteListAlerts(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteListAlerts", reflect.TypeOf((*MockIndexer)(nil).DeleteListAlerts), ids)
}

// Search mocks base method.
func (m *MockIndexer) Search(ctx context.Context, q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, q}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Search", varargs...)
	ret0, _ := ret[0].([]search.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockIndexerMockRecorder) Search(ctx, q interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, q}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIndexer)(nil).Search), varargs...)
}

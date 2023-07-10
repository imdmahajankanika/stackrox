// Code generated by MockGen. DO NOT EDIT.
// Source: searcher.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRanker is a mock of Ranker interface.
type MockRanker struct {
	ctrl     *gomock.Controller
	recorder *MockRankerMockRecorder
}

// MockRankerMockRecorder is the mock recorder for MockRanker.
type MockRankerMockRecorder struct {
	mock *MockRanker
}

// NewMockRanker creates a new mock instance.
func NewMockRanker(ctrl *gomock.Controller) *MockRanker {
	mock := &MockRanker{ctrl: ctrl}
	mock.recorder = &MockRankerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRanker) EXPECT() *MockRankerMockRecorder {
	return m.recorder
}

// GetRankForID mocks base method.
func (m *MockRanker) GetRankForID(from string) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRankForID", from)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetRankForID indicates an expected call of GetRankForID.
func (mr *MockRankerMockRecorder) GetRankForID(from interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRankForID", reflect.TypeOf((*MockRanker)(nil).GetRankForID), from)
}

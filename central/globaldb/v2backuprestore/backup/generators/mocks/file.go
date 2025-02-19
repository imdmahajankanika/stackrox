// Code generated by MockGen. DO NOT EDIT.
// Source: file.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFileGenerator is a mock of FileGenerator interface.
type MockFileGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockFileGeneratorMockRecorder
}

// MockFileGeneratorMockRecorder is the mock recorder for MockFileGenerator.
type MockFileGeneratorMockRecorder struct {
	mock *MockFileGenerator
}

// NewMockFileGenerator creates a new mock instance.
func NewMockFileGenerator(ctrl *gomock.Controller) *MockFileGenerator {
	mock := &MockFileGenerator{ctrl: ctrl}
	mock.recorder = &MockFileGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileGenerator) EXPECT() *MockFileGeneratorMockRecorder {
	return m.recorder
}

// WriteFile mocks base method.
func (m *MockFileGenerator) WriteFile(ctx context.Context, path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFile", ctx, path)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFile indicates an expected call of WriteFile.
func (mr *MockFileGeneratorMockRecorder) WriteFile(ctx, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockFileGenerator)(nil).WriteFile), ctx, path)
}

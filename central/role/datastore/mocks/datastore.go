// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/central/role/datastore (interfaces: DataStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	storage "github.com/stackrox/rox/generated/storage"
	reflect "reflect"
)

// MockDataStore is a mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// AddRole mocks base method
func (m *MockDataStore) AddRole(arg0 context.Context, arg1 *storage.Role) error {
	ret := m.ctrl.Call(m, "AddRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRole indicates an expected call of AddRole
func (mr *MockDataStoreMockRecorder) AddRole(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRole", reflect.TypeOf((*MockDataStore)(nil).AddRole), arg0, arg1)
}

// GetAllRoles mocks base method
func (m *MockDataStore) GetAllRoles(arg0 context.Context) ([]*storage.Role, error) {
	ret := m.ctrl.Call(m, "GetAllRoles", arg0)
	ret0, _ := ret[0].([]*storage.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRoles indicates an expected call of GetAllRoles
func (mr *MockDataStoreMockRecorder) GetAllRoles(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRoles", reflect.TypeOf((*MockDataStore)(nil).GetAllRoles), arg0)
}

// GetRole mocks base method
func (m *MockDataStore) GetRole(arg0 context.Context, arg1 string) (*storage.Role, error) {
	ret := m.ctrl.Call(m, "GetRole", arg0, arg1)
	ret0, _ := ret[0].(*storage.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole
func (mr *MockDataStoreMockRecorder) GetRole(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockDataStore)(nil).GetRole), arg0, arg1)
}

// RemoveRole mocks base method
func (m *MockDataStore) RemoveRole(arg0 context.Context, arg1 string) error {
	ret := m.ctrl.Call(m, "RemoveRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRole indicates an expected call of RemoveRole
func (mr *MockDataStoreMockRecorder) RemoveRole(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRole", reflect.TypeOf((*MockDataStore)(nil).RemoveRole), arg0, arg1)
}

// UpdateRole mocks base method
func (m *MockDataStore) UpdateRole(arg0 context.Context, arg1 *storage.Role) error {
	ret := m.ctrl.Call(m, "UpdateRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRole indicates an expected call of UpdateRole
func (mr *MockDataStoreMockRecorder) UpdateRole(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRole", reflect.TypeOf((*MockDataStore)(nil).UpdateRole), arg0, arg1)
}

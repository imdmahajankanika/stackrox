// Code generated by MockGen. DO NOT EDIT.
// Source: crud.go

// Package mocks is a generated GoMock package.
package mocks

import (
	proto "github.com/gogo/protobuf/proto"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCrud is a mock of Crud interface
type MockCrud struct {
	ctrl     *gomock.Controller
	recorder *MockCrudMockRecorder
}

// MockCrudMockRecorder is the mock recorder for MockCrud
type MockCrudMockRecorder struct {
	mock *MockCrud
}

// NewMockCrud creates a new mock instance
func NewMockCrud(ctrl *gomock.Controller) *MockCrud {
	mock := &MockCrud{ctrl: ctrl}
	mock.recorder = &MockCrudMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCrud) EXPECT() *MockCrudMockRecorder {
	return m.recorder
}

// Count mocks base method
func (m *MockCrud) Count() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockCrudMockRecorder) Count() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockCrud)(nil).Count))
}

// Exists mocks base method
func (m *MockCrud) Exists(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists
func (mr *MockCrudMockRecorder) Exists(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockCrud)(nil).Exists), id)
}

// GetKeys mocks base method
func (m *MockCrud) GetKeys() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeys")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeys indicates an expected call of GetKeys
func (mr *MockCrudMockRecorder) GetKeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeys", reflect.TypeOf((*MockCrud)(nil).GetKeys))
}

// Get mocks base method
func (m *MockCrud) Get(id string) (proto.Message, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(proto.Message)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get
func (mr *MockCrudMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCrud)(nil).Get), id)
}

// GetMany mocks base method
func (m *MockCrud) GetMany(ids []string) ([]proto.Message, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMany", ids)
	ret0, _ := ret[0].([]proto.Message)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMany indicates an expected call of GetMany
func (mr *MockCrudMockRecorder) GetMany(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMany", reflect.TypeOf((*MockCrud)(nil).GetMany), ids)
}

// Walk mocks base method
func (m *MockCrud) Walk(arg0 func(proto.Message) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Walk", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Walk indicates an expected call of Walk
func (mr *MockCrudMockRecorder) Walk(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Walk", reflect.TypeOf((*MockCrud)(nil).Walk), arg0)
}

// Upsert mocks base method
func (m *MockCrud) Upsert(kv proto.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", kv)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert
func (mr *MockCrudMockRecorder) Upsert(kv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockCrud)(nil).Upsert), kv)
}

// UpsertMany mocks base method
func (m *MockCrud) UpsertMany(msgs []proto.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertMany", msgs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertMany indicates an expected call of UpsertMany
func (mr *MockCrudMockRecorder) UpsertMany(msgs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertMany", reflect.TypeOf((*MockCrud)(nil).UpsertMany), msgs)
}

// UpsertWithID mocks base method
func (m *MockCrud) UpsertWithID(id string, msg proto.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertWithID", id, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertWithID indicates an expected call of UpsertWithID
func (mr *MockCrudMockRecorder) UpsertWithID(id, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertWithID", reflect.TypeOf((*MockCrud)(nil).UpsertWithID), id, msg)
}

// UpsertManyWithIDs mocks base method
func (m *MockCrud) UpsertManyWithIDs(ids []string, msgs []proto.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertManyWithIDs", ids, msgs)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertManyWithIDs indicates an expected call of UpsertManyWithIDs
func (mr *MockCrudMockRecorder) UpsertManyWithIDs(ids, msgs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertManyWithIDs", reflect.TypeOf((*MockCrud)(nil).UpsertManyWithIDs), ids, msgs)
}

// Delete mocks base method
func (m *MockCrud) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockCrudMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCrud)(nil).Delete), id)
}

// DeleteMany mocks base method
func (m *MockCrud) DeleteMany(ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMany", ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMany indicates an expected call of DeleteMany
func (mr *MockCrudMockRecorder) DeleteMany(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMany", reflect.TypeOf((*MockCrud)(nil).DeleteMany), ids)
}

// AckKeysIndexed mocks base method
func (m *MockCrud) AckKeysIndexed(keys ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AckKeysIndexed", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AckKeysIndexed indicates an expected call of AckKeysIndexed
func (mr *MockCrudMockRecorder) AckKeysIndexed(keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AckKeysIndexed", reflect.TypeOf((*MockCrud)(nil).AckKeysIndexed), keys...)
}

// GetKeysToIndex mocks base method
func (m *MockCrud) GetKeysToIndex() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeysToIndex")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeysToIndex indicates an expected call of GetKeysToIndex
func (mr *MockCrudMockRecorder) GetKeysToIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeysToIndex", reflect.TypeOf((*MockCrud)(nil).GetKeysToIndex))
}

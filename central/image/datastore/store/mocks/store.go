// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	storage "github.com/stackrox/rox/generated/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockStore) Count(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockStoreMockRecorder) Count(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockStore)(nil).Count), ctx)
}

// Delete mocks base method.
func (m *MockStore) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStoreMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStore)(nil).Delete), ctx, id)
}

// Exists mocks base method.
func (m *MockStore) Exists(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockStoreMockRecorder) Exists(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockStore)(nil).Exists), ctx, id)
}

// Get mocks base method.
func (m *MockStore) Get(ctx context.Context, id string) (*storage.Image, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*storage.Image)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), ctx, id)
}

// GetImageMetadata mocks base method.
func (m *MockStore) GetImageMetadata(ctx context.Context, id string) (*storage.Image, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageMetadata", ctx, id)
	ret0, _ := ret[0].(*storage.Image)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetImageMetadata indicates an expected call of GetImageMetadata.
func (mr *MockStoreMockRecorder) GetImageMetadata(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageMetadata", reflect.TypeOf((*MockStore)(nil).GetImageMetadata), ctx, id)
}

// GetMany mocks base method.
func (m *MockStore) GetMany(ctx context.Context, ids []string) ([]*storage.Image, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMany", ctx, ids)
	ret0, _ := ret[0].([]*storage.Image)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMany indicates an expected call of GetMany.
func (mr *MockStoreMockRecorder) GetMany(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMany", reflect.TypeOf((*MockStore)(nil).GetMany), ctx, ids)
}

// GetManyImageMetadata mocks base method.
func (m *MockStore) GetManyImageMetadata(ctx context.Context, id []string) ([]*storage.Image, []int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManyImageMetadata", ctx, id)
	ret0, _ := ret[0].([]*storage.Image)
	ret1, _ := ret[1].([]int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetManyImageMetadata indicates an expected call of GetManyImageMetadata.
func (mr *MockStoreMockRecorder) GetManyImageMetadata(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManyImageMetadata", reflect.TypeOf((*MockStore)(nil).GetManyImageMetadata), ctx, id)
}

// UpdateVulnState mocks base method.
func (m *MockStore) UpdateVulnState(ctx context.Context, cve string, imageIDs []string, state storage.VulnerabilityState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVulnState", ctx, cve, imageIDs, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVulnState indicates an expected call of UpdateVulnState.
func (mr *MockStoreMockRecorder) UpdateVulnState(ctx, cve, imageIDs, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVulnState", reflect.TypeOf((*MockStore)(nil).UpdateVulnState), ctx, cve, imageIDs, state)
}

// Upsert mocks base method.
func (m *MockStore) Upsert(ctx context.Context, image *storage.Image) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", ctx, image)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert.
func (mr *MockStoreMockRecorder) Upsert(ctx, image interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockStore)(nil).Upsert), ctx, image)
}

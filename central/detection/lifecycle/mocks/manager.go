// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	storage "github.com/stackrox/rox/generated/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// DeploymentRemoved mocks base method.
func (m *MockManager) DeploymentRemoved(deploymentID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeploymentRemoved", deploymentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeploymentRemoved indicates an expected call of DeploymentRemoved.
func (mr *MockManagerMockRecorder) DeploymentRemoved(deploymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeploymentRemoved", reflect.TypeOf((*MockManager)(nil).DeploymentRemoved), deploymentID)
}

// HandleDeploymentAlerts mocks base method.
func (m *MockManager) HandleDeploymentAlerts(deploymentID string, alerts []*storage.Alert, stage storage.LifecycleStage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleDeploymentAlerts", deploymentID, alerts, stage)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleDeploymentAlerts indicates an expected call of HandleDeploymentAlerts.
func (mr *MockManagerMockRecorder) HandleDeploymentAlerts(deploymentID, alerts, stage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleDeploymentAlerts", reflect.TypeOf((*MockManager)(nil).HandleDeploymentAlerts), deploymentID, alerts, stage)
}

// HandleResourceAlerts mocks base method.
func (m *MockManager) HandleResourceAlerts(clusterID string, alerts []*storage.Alert, stage storage.LifecycleStage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleResourceAlerts", clusterID, alerts, stage)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleResourceAlerts indicates an expected call of HandleResourceAlerts.
func (mr *MockManagerMockRecorder) HandleResourceAlerts(clusterID, alerts, stage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleResourceAlerts", reflect.TypeOf((*MockManager)(nil).HandleResourceAlerts), clusterID, alerts, stage)
}

// IndicatorAdded mocks base method.
func (m *MockManager) IndicatorAdded(indicator *storage.ProcessIndicator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndicatorAdded", indicator)
	ret0, _ := ret[0].(error)
	return ret0
}

// IndicatorAdded indicates an expected call of IndicatorAdded.
func (mr *MockManagerMockRecorder) IndicatorAdded(indicator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndicatorAdded", reflect.TypeOf((*MockManager)(nil).IndicatorAdded), indicator)
}

// RemoveDeploymentFromObservation mocks base method.
func (m *MockManager) RemoveDeploymentFromObservation(deploymentID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveDeploymentFromObservation", deploymentID)
}

// RemoveDeploymentFromObservation indicates an expected call of RemoveDeploymentFromObservation.
func (mr *MockManagerMockRecorder) RemoveDeploymentFromObservation(deploymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDeploymentFromObservation", reflect.TypeOf((*MockManager)(nil).RemoveDeploymentFromObservation), deploymentID)
}

// RemovePolicy mocks base method.
func (m *MockManager) RemovePolicy(policyID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePolicy", policyID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePolicy indicates an expected call of RemovePolicy.
func (mr *MockManagerMockRecorder) RemovePolicy(policyID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePolicy", reflect.TypeOf((*MockManager)(nil).RemovePolicy), policyID)
}

// UpsertPolicy mocks base method.
func (m *MockManager) UpsertPolicy(policy *storage.Policy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertPolicy", policy)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertPolicy indicates an expected call of UpsertPolicy.
func (mr *MockManagerMockRecorder) UpsertPolicy(policy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertPolicy", reflect.TypeOf((*MockManager)(nil).UpsertPolicy), policy)
}

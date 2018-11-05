// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

import timestamp "github.com/stackrox/rox/pkg/timestamp"
import types "github.com/gogo/protobuf/types"
import v1 "github.com/stackrox/rox/generated/api/v1"

// FlowStore is an autogenerated mock type for the FlowStore type
type FlowStore struct {
	mock.Mock
}

// GetAllFlows provides a mock function with given fields:
func (_m *FlowStore) GetAllFlows() ([]*v1.NetworkFlow, types.Timestamp, error) {
	ret := _m.Called()

	var r0 []*v1.NetworkFlow
	if rf, ok := ret.Get(0).(func() []*v1.NetworkFlow); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.NetworkFlow)
		}
	}

	var r1 types.Timestamp
	if rf, ok := ret.Get(1).(func() types.Timestamp); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(types.Timestamp)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetFlow provides a mock function with given fields: props
func (_m *FlowStore) GetFlow(props *v1.NetworkFlowProperties) (*v1.NetworkFlow, error) {
	ret := _m.Called(props)

	var r0 *v1.NetworkFlow
	if rf, ok := ret.Get(0).(func(*v1.NetworkFlowProperties) *v1.NetworkFlow); ok {
		r0 = rf(props)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.NetworkFlow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.NetworkFlowProperties) error); ok {
		r1 = rf(props)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFlow provides a mock function with given fields: props
func (_m *FlowStore) RemoveFlow(props *v1.NetworkFlowProperties) error {
	ret := _m.Called(props)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.NetworkFlowProperties) error); ok {
		r0 = rf(props)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertFlows provides a mock function with given fields: flows, lastUpdateTS
func (_m *FlowStore) UpsertFlows(flows []*v1.NetworkFlow, lastUpdateTS timestamp.MicroTS) error {
	ret := _m.Called(flows, lastUpdateTS)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*v1.NetworkFlow, timestamp.MicroTS) error); ok {
		r0 = rf(flows, lastUpdateTS)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import dnrintegration "github.com/stackrox/rox/central/dnrintegration"
import mock "github.com/stretchr/testify/mock"
import v1 "github.com/stackrox/rox/generated/api/v1"

// DataStore is an autogenerated mock type for the DataStore type
type DataStore struct {
	mock.Mock
}

// AddDNRIntegration provides a mock function with given fields: proto, integration
func (_m *DataStore) AddDNRIntegration(proto *v1.DNRIntegration, integration dnrintegration.DNRIntegration) (string, error) {
	ret := _m.Called(proto, integration)

	var r0 string
	if rf, ok := ret.Get(0).(func(*v1.DNRIntegration, dnrintegration.DNRIntegration) string); ok {
		r0 = rf(proto, integration)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.DNRIntegration, dnrintegration.DNRIntegration) error); ok {
		r1 = rf(proto, integration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForCluster provides a mock function with given fields: clusterID
func (_m *DataStore) ForCluster(clusterID string) (dnrintegration.DNRIntegration, bool) {
	ret := _m.Called(clusterID)

	var r0 dnrintegration.DNRIntegration
	if rf, ok := ret.Get(0).(func(string) dnrintegration.DNRIntegration); ok {
		r0 = rf(clusterID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(dnrintegration.DNRIntegration)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(clusterID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetDNRIntegration provides a mock function with given fields: id
func (_m *DataStore) GetDNRIntegration(id string) (*v1.DNRIntegration, bool, error) {
	ret := _m.Called(id)

	var r0 *v1.DNRIntegration
	if rf, ok := ret.Get(0).(func(string) *v1.DNRIntegration); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.DNRIntegration)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetDNRIntegrations provides a mock function with given fields: request
func (_m *DataStore) GetDNRIntegrations(request *v1.GetDNRIntegrationsRequest) ([]*v1.DNRIntegration, error) {
	ret := _m.Called(request)

	var r0 []*v1.DNRIntegration
	if rf, ok := ret.Get(0).(func(*v1.GetDNRIntegrationsRequest) []*v1.DNRIntegration); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.DNRIntegration)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.GetDNRIntegrationsRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveDNRIntegration provides a mock function with given fields: id
func (_m *DataStore) RemoveDNRIntegration(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateDNRIntegration provides a mock function with given fields: proto, integration
func (_m *DataStore) UpdateDNRIntegration(proto *v1.DNRIntegration, integration dnrintegration.DNRIntegration) error {
	ret := _m.Called(proto, integration)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.DNRIntegration, dnrintegration.DNRIntegration) error); ok {
		r0 = rf(proto, integration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

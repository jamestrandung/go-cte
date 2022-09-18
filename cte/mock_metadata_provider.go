// Code generated by mockery v2.14.0. DO NOT EDIT.

package cte

import mock "github.com/stretchr/testify/mock"

// MockMetadataProvider is an autogenerated mock type for the MetadataProvider type
type MockMetadataProvider struct {
	mock.Mock
}

// CTEMetadata provides a mock function with given fields:
func (_m *MockMetadataProvider) CTEMetadata() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

type mockConstructorTestingTNewMockMetadataProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockMetadataProvider creates a new instance of MockMetadataProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockMetadataProvider(t mockConstructorTestingTNewMockMetadataProvider) *MockMetadataProvider {
	mock := &MockMetadataProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

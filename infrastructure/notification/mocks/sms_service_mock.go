// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	notification "github.com/machtwatch/catalyst-go-skeleton/infrastructure/notification"
	mock "github.com/stretchr/testify/mock"
)

// ISMSService is an autogenerated mock type for the ISMSService type
type ISMSService struct {
	mock.Mock
}

// Send provides a mock function with given fields: request
func (_m *ISMSService) Send(request notification.SMSRequest) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(notification.SMSRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewISMSService interface {
	mock.TestingT
	Cleanup(func())
}

// NewISMSService creates a new instance of ISMSService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewISMSService(t mockConstructorTestingTNewISMSService) *ISMSService {
	mock := &ISMSService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

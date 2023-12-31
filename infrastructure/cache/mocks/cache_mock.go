// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Cache is an autogenerated mock type for the Cache type
type Cache struct {
	mock.Mock
}

// DEL provides a mock function with given fields: ctx, key
func (_m *Cache) DEL(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GET provides a mock function with given fields: ctx, key
func (_m *Cache) GET(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PING provides a mock function with given fields: ctx
func (_m *Cache) PING(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SETEX provides a mock function with given fields: ctx, key, value, expiration
func (_m *Cache) SETEX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ret := _m.Called(ctx, key, value, expiration)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) error); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCache interface {
	mock.TestingT
	Cleanup(func())
}

// NewCache creates a new instance of Cache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCache(t mockConstructorTestingTNewCache) *Cache {
	mock := &Cache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	uow "github.com/rafael-piovesan/go-rocket-ride/v2/datastore/uow"
)

// UnitOfWork is an autogenerated mock type for the UnitOfWork type
type UnitOfWork struct {
	mock.Mock
}

// Do provides a mock function with given fields: _a0, _a1
func (_m *UnitOfWork) Do(_a0 context.Context, _a1 uow.UnitOfWorkBlock) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uow.UnitOfWorkBlock) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUnitOfWork creates a new instance of UnitOfWork. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewUnitOfWork(t testing.TB) *UnitOfWork {
	mock := &UnitOfWork{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

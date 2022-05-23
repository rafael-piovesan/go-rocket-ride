// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/rafael-piovesan/go-rocket-ride/v2/entity"

	mock "github.com/stretchr/testify/mock"

	repo "github.com/rafael-piovesan/go-rocket-ride/v2/pkg/repo"

	testing "testing"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *User) Delete(_a0 context.Context, _a1 *entity.User) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: _a0, _a1
func (_m *User) FindAll(_a0 context.Context, _a1 ...repo.SelectCriteria) ([]entity.User, error) {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []entity.User
	if rf, ok := ret.Get(0).(func(context.Context, ...repo.SelectCriteria) []entity.User); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...repo.SelectCriteria) error); ok {
		r1 = rf(_a0, _a1...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: _a0, _a1
func (_m *User) FindOne(_a0 context.Context, _a1 ...repo.SelectCriteria) (entity.User, error) {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(context.Context, ...repo.SelectCriteria) entity.User); ok {
		r0 = rf(_a0, _a1...)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...repo.SelectCriteria) error); ok {
		r1 = rf(_a0, _a1...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *User) Save(_a0 context.Context, _a1 *entity.User) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *User) Update(_a0 context.Context, _a1 *entity.User) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUser creates a new instance of User. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewUser(t testing.TB) *User {
	mock := &User{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
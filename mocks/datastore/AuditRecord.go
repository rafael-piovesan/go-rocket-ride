// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"

	data "github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"

	entity "github.com/rafael-piovesan/go-rocket-ride/v2/entity"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// AuditRecord is an autogenerated mock type for the AuditRecord type
type AuditRecord struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *AuditRecord) Delete(_a0 context.Context, _a1 *entity.AuditRecord) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.AuditRecord) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: _a0, _a1
func (_m *AuditRecord) FindAll(_a0 context.Context, _a1 ...data.SelectCriteria) ([]entity.AuditRecord, error) {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []entity.AuditRecord
	if rf, ok := ret.Get(0).(func(context.Context, ...data.SelectCriteria) []entity.AuditRecord); ok {
		r0 = rf(_a0, _a1...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.AuditRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...data.SelectCriteria) error); ok {
		r1 = rf(_a0, _a1...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: _a0, _a1
func (_m *AuditRecord) FindOne(_a0 context.Context, _a1 ...data.SelectCriteria) (entity.AuditRecord, error) {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 entity.AuditRecord
	if rf, ok := ret.Get(0).(func(context.Context, ...data.SelectCriteria) entity.AuditRecord); ok {
		r0 = rf(_a0, _a1...)
	} else {
		r0 = ret.Get(0).(entity.AuditRecord)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...data.SelectCriteria) error); ok {
		r1 = rf(_a0, _a1...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *AuditRecord) Save(_a0 context.Context, _a1 *entity.AuditRecord) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.AuditRecord) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *AuditRecord) Update(_a0 context.Context, _a1 *entity.AuditRecord) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.AuditRecord) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAuditRecord creates a new instance of AuditRecord. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuditRecord(t testing.TB) *AuditRecord {
	mock := &AuditRecord{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

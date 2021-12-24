// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/rafael-piovesan/go-rocket-ride/entity"
	mock "github.com/stretchr/testify/mock"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
)

// Datastore is an autogenerated mock type for the Datastore type
type Datastore struct {
	mock.Mock
}

// Atomic provides a mock function with given fields: ctx, fn
func (_m *Datastore) Atomic(ctx context.Context, fn func(rocketride.Datastore) error) error {
	ret := _m.Called(ctx, fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(rocketride.Datastore) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAuditRecord provides a mock function with given fields: ctx, ar
func (_m *Datastore) CreateAuditRecord(ctx context.Context, ar *entity.AuditRecord) (*entity.AuditRecord, error) {
	ret := _m.Called(ctx, ar)

	var r0 *entity.AuditRecord
	if rf, ok := ret.Get(0).(func(context.Context, *entity.AuditRecord) *entity.AuditRecord); ok {
		r0 = rf(ctx, ar)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.AuditRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.AuditRecord) error); ok {
		r1 = rf(ctx, ar)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateIdempotencyKey provides a mock function with given fields: ctx, ik
func (_m *Datastore) CreateIdempotencyKey(ctx context.Context, ik *entity.IdempotencyKey) (*entity.IdempotencyKey, error) {
	ret := _m.Called(ctx, ik)

	var r0 *entity.IdempotencyKey
	if rf, ok := ret.Get(0).(func(context.Context, *entity.IdempotencyKey) *entity.IdempotencyKey); ok {
		r0 = rf(ctx, ik)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.IdempotencyKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.IdempotencyKey) error); ok {
		r1 = rf(ctx, ik)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRide provides a mock function with given fields: ctx, rd
func (_m *Datastore) CreateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	ret := _m.Called(ctx, rd)

	var r0 *entity.Ride
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Ride) *entity.Ride); ok {
		r0 = rf(ctx, rd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Ride)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.Ride) error); ok {
		r1 = rf(ctx, rd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateStagedJob provides a mock function with given fields: ctx, sj
func (_m *Datastore) CreateStagedJob(ctx context.Context, sj *entity.StagedJob) (*entity.StagedJob, error) {
	ret := _m.Called(ctx, sj)

	var r0 *entity.StagedJob
	if rf, ok := ret.Get(0).(func(context.Context, *entity.StagedJob) *entity.StagedJob); ok {
		r0 = rf(ctx, sj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.StagedJob)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.StagedJob) error); ok {
		r1 = rf(ctx, sj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetIdempotencyKey provides a mock function with given fields: ctx, key, userID
func (_m *Datastore) GetIdempotencyKey(ctx context.Context, key string, userID int64) (*entity.IdempotencyKey, error) {
	ret := _m.Called(ctx, key, userID)

	var r0 *entity.IdempotencyKey
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) *entity.IdempotencyKey); ok {
		r0 = rf(ctx, key, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.IdempotencyKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, key, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRideByIdempotencyKeyID provides a mock function with given fields: ctx, keyID
func (_m *Datastore) GetRideByIdempotencyKeyID(ctx context.Context, keyID int64) (*entity.Ride, error) {
	ret := _m.Called(ctx, keyID)

	var r0 *entity.Ride
	if rf, ok := ret.Get(0).(func(context.Context, int64) *entity.Ride); ok {
		r0 = rf(ctx, keyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Ride)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, keyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: ctx, email
func (_m *Datastore) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateIdempotencyKey provides a mock function with given fields: ctx, ik
func (_m *Datastore) UpdateIdempotencyKey(ctx context.Context, ik *entity.IdempotencyKey) (*entity.IdempotencyKey, error) {
	ret := _m.Called(ctx, ik)

	var r0 *entity.IdempotencyKey
	if rf, ok := ret.Get(0).(func(context.Context, *entity.IdempotencyKey) *entity.IdempotencyKey); ok {
		r0 = rf(ctx, ik)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.IdempotencyKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.IdempotencyKey) error); ok {
		r1 = rf(ctx, ik)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRide provides a mock function with given fields: ctx, rd
func (_m *Datastore) UpdateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	ret := _m.Called(ctx, rd)

	var r0 *entity.Ride
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Ride) *entity.Ride); ok {
		r0 = rf(ctx, rd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Ride)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.Ride) error); ok {
		r1 = rf(ctx, rd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

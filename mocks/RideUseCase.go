// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/rafael-piovesan/go-rocket-ride/entity"
	mock "github.com/stretchr/testify/mock"
)

// RideUseCase is an autogenerated mock type for the RideUseCase type
type RideUseCase struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, ik, rd
func (_m *RideUseCase) Create(ctx context.Context, ik *entity.IdempotencyKey, rd *entity.Ride) (*entity.IdempotencyKey, error) {
	ret := _m.Called(ctx, ik, rd)

	var r0 *entity.IdempotencyKey
	if rf, ok := ret.Get(0).(func(context.Context, *entity.IdempotencyKey, *entity.Ride) *entity.IdempotencyKey); ok {
		r0 = rf(ctx, ik, rd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.IdempotencyKey)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.IdempotencyKey, *entity.Ride) error); ok {
		r1 = rf(ctx, ik, rd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
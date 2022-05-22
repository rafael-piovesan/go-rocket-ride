package rocketride

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

type Datastore interface {
	Atomic(ctx context.Context, fn func(ds Datastore) error) error
	CreateAuditRecord(ctx context.Context, ar *entity.AuditRecord) (*entity.AuditRecord, error)
	CreateIdempotencyKey(ctx context.Context, ik *entity.IdempotencyKey) (*entity.IdempotencyKey, error)
	CreateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error)
	CreateStagedJob(ctx context.Context, sj *entity.StagedJob) (*entity.StagedJob, error)
	GetIdempotencyKey(ctx context.Context, key string, userID int64) (*entity.IdempotencyKey, error)
	GetRideByIdempotencyKeyID(ctx context.Context, keyID int64) (*entity.Ride, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateIdempotencyKey(ctx context.Context, ik *entity.IdempotencyKey) (*entity.IdempotencyKey, error)
	UpdateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error)
}

type RideUseCase interface {
	Create(ctx context.Context, ik *entity.IdempotencyKey, rd *entity.Ride) (*entity.IdempotencyKey, error)
}

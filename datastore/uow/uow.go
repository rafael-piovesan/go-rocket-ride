package uow

import (
	"context"
	"database/sql"

	"github.com/rafael-piovesan/go-rocket-ride/v2/datastore"
	"github.com/uptrace/bun"
)

type uowStore struct {
	auditRecords datastore.AuditRecord
	idemKeys     datastore.IdempotencyKey
	rides        datastore.Ride
	stagedJobs   datastore.StagedJob
	users        datastore.User
}

type UOWStore interface {
	AuditRecords() datastore.AuditRecord
	IdempotencyKeys() datastore.IdempotencyKey
	Rides() datastore.Ride
	StagedJobs() datastore.StagedJob
	Users() datastore.User
}

func (u uowStore) AuditRecords() datastore.AuditRecord {
	return u.auditRecords
}

func (u uowStore) IdempotencyKeys() datastore.IdempotencyKey {
	return u.idemKeys
}

func (u uowStore) Rides() datastore.Ride {
	return u.rides
}

func (u uowStore) StagedJobs() datastore.StagedJob {
	return u.stagedJobs
}

func (u uowStore) Users() datastore.User {
	return u.users
}

type UOWBlock func(store UOWStore) error

type unitOfWork struct {
	conn *bun.DB
	db   bun.IDB
}

type UnitOfWork interface {
	Do(context.Context, UOWBlock) error
}

func New(db *bun.DB) UnitOfWork {
	return &unitOfWork{
		conn: db,
		db:   db,
	}
}

func (s *unitOfWork) Do(ctx context.Context, fn UOWBlock) error {
	return s.conn.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		newStore := &uowStore{
			auditRecords: datastore.NewAuditRecord(tx),
			idemKeys:     datastore.NewIdempotencyKey(tx),
			rides:        datastore.NewRide(tx),
			stagedJobs:   datastore.NewStagedJob(tx),
			users:        datastore.NewUser(tx),
		}
		return fn(newStore)
	})
}

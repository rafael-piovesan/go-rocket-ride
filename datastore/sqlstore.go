package datastore

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type atomicStore struct {
	auditRecords AuditRecord
	idemKeys     IdempotencyKey
	rides        Ride
	stagedJobs   StagedJob
	users        User
}

type AtomicStore interface {
	AuditRecords() AuditRecord
	IdempotencyKeys() IdempotencyKey
	Rides() Ride
	StagedJobs() StagedJob
	Users() User
}

func (a *atomicStore) AuditRecords() AuditRecord {
	return a.auditRecords
}

func (a *atomicStore) IdempotencyKeys() IdempotencyKey {
	return a.idemKeys
}

func (a *atomicStore) Rides() Ride {
	return a.rides
}

func (a *atomicStore) StagedJobs() StagedJob {
	return a.stagedJobs
}

func (a *atomicStore) Users() User {
	return a.users
}

type AtomicBlock func(store AtomicStore) error

type sqlStore struct {
	conn *bun.DB
	db   bun.IDB
}

type Store interface {
	Atomic(context.Context, AtomicBlock) error
}

func New(db *bun.DB) Store {
	return &sqlStore{
		conn: db,
		db:   db,
	}
}

func (s *sqlStore) Atomic(ctx context.Context, fn AtomicBlock) error {
	return s.conn.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		newStore := &atomicStore{
			auditRecords: NewAuditRecord(tx),
			idemKeys:     NewIdempotencyKey(tx),
			rides:        NewRide(tx),
			stagedJobs:   NewStagedJob(tx),
			users:        NewUser(tx),
		}
		return fn(newStore)
	})
}

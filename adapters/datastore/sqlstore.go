package datastore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
)

type sqlStore struct {
	conn *bun.DB
	db   bun.IDB
}

func NewStore(db *bun.DB) rocketride.Datastore {
	return &sqlStore{
		conn: db,
		db:   db,
	}
}

func (s *sqlStore) Atomic(ctx context.Context, fn func(store rocketride.Datastore) error) (err error) {
	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic err: %v", p)
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	// TODO: check if it works for nested transactions as well
	newStore := &sqlStore{
		conn: s.conn,
		db:   tx,
	}
	err = fn(newStore)
	return err
}

var _ rocketride.Datastore = (*sqlStore)(nil)

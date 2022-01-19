package bun

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	// loading bun's official Postgres driver.
	_ "github.com/uptrace/bun/driver/pgdriver"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"
)

type sqlStore struct {
	conn *bun.DB
	db   bun.IDB
}

func NewStore(dsn string) (s rocketride.Datastore, err error) {
	sqldb, err := sql.Open("pg", dsn)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	s = &sqlStore{
		conn: db,
		db:   db,
	}
	return
}

func (s *sqlStore) Atomic(ctx context.Context, fn func(store rocketride.Datastore) error) (err error) {
	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()

			switch e := p.(type) {
			case runtime.Error:
				panic(e)
			case error:
				err = fmt.Errorf("panic err: %v", p)
				return
			default:
				panic(e)
			}
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

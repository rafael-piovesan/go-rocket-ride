package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	rocketride "github.com/rafael-piovesan/go-rocket-ride"

	// loading Postgres driver.
	_ "github.com/lib/pq"
)

type sqlStore struct {
	conn *sql.DB
	q    *Queries
}

func NewStore(dsn string) (s rocketride.Datastore, err error) {
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	s = &sqlStore{
		conn: sqldb,
		q:    New(sqldb),
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
		q:    New(tx),
	}
	err = fn(newStore)
	return err
}

var _ rocketride.Datastore = (*sqlStore)(nil)

package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

var ErrRecordNotFound = errors.New("record not found")

type SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery

type CRUDRepository[T any] struct {
	DB bun.IDB
}

func New[T any](db bun.IDB) CRUDRepository[T] {
	return CRUDRepository[T]{DB: db}
}

func (c CRUDRepository[T]) FindAll(ctx context.Context, sc ...SelectCriteria) ([]T, error) {
	var rows []T

	q := c.DB.NewSelect().Model(&rows)
	for i := range sc {
		q.Apply(sc[i])
	}

	err := q.Scan(ctx)
	return rows, err
}

func (c CRUDRepository[T]) FindOne(ctx context.Context, sc ...SelectCriteria) (T, error) {
	var row T

	q := c.DB.NewSelect().Model(&row)
	for i := range sc {
		q.Apply(sc[i])
	}

	err := q.Limit(1).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return row, ErrRecordNotFound
	}
	return row, err
}

func (c CRUDRepository[T]) Save(ctx context.Context, model *T) error {
	_, err := c.DB.NewInsert().Model(model).Returning("*").Exec(ctx)
	return err
}

func (c CRUDRepository[T]) Delete(ctx context.Context, model *T) error {
	_, err := c.DB.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

func (c CRUDRepository[T]) Update(ctx context.Context, model *T) error {
	_, err := c.DB.NewUpdate().Model(model).WherePK().Returning("*").Exec(ctx)
	return err
}

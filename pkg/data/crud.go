package data

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

var ErrRecordNotFound = errors.New("record not found")

type SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery

type ICRUDStore[T any] interface {
	FindAll(context.Context, ...SelectCriteria) ([]T, error)
	FindOne(context.Context, ...SelectCriteria) (T, error)
	Delete(context.Context, *T) error
	Save(context.Context, *T) error
	Update(context.Context, *T) error
}

type CRUDStore[T any] struct {
	DB bun.IDB
}

func New[T any](db bun.IDB) ICRUDStore[T] {
	return CRUDStore[T]{DB: db}
}

func (c CRUDStore[T]) FindAll(ctx context.Context, sc ...SelectCriteria) ([]T, error) {
	var rows []T

	q := c.DB.NewSelect().Model(&rows)
	for i := range sc {
		q.Apply(sc[i])
	}

	err := q.Scan(ctx)
	return rows, err
}

func (c CRUDStore[T]) FindOne(ctx context.Context, sc ...SelectCriteria) (T, error) {
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

func (c CRUDStore[T]) Save(ctx context.Context, model *T) error {
	_, err := c.DB.NewInsert().Model(model).Returning("*").Exec(ctx)
	return err
}

func (c CRUDStore[T]) Delete(ctx context.Context, model *T) error {
	_, err := c.DB.NewDelete().Model(model).WherePK().Exec(ctx)
	return err
}

func (c CRUDStore[T]) Update(ctx context.Context, model *T) error {
	_, err := c.DB.NewUpdate().Model(model).WherePK().Returning("*").Exec(ctx)
	return err
}

package datastore

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/repo"
	"github.com/uptrace/bun"
)

type IdempotencyKey interface {
	FindAll(context.Context, ...repo.SelectCriteria) ([]entity.IdempotencyKey, error)
	FindOne(context.Context, ...repo.SelectCriteria) (entity.IdempotencyKey, error)
	Delete(context.Context, *entity.IdempotencyKey) error
	Save(context.Context, *entity.IdempotencyKey) error
	Update(context.Context, *entity.IdempotencyKey) error
}

type idempotencyKeyStore struct {
	*repo.CRUDRepository[entity.IdempotencyKey]
}

func NewIdempotencyKey(db bun.IDB) IdempotencyKey {
	return idempotencyKeyStore{&repo.CRUDRepository[entity.IdempotencyKey]{DB: db}}
}

func IdemKeyWithKey(key string) repo.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("idempotency_key = ?", key)
	}
}

func IdemKeyWithUserID(uid int64) repo.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("user_id = ?", uid)
	}
}

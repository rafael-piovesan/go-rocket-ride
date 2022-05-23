package datastore

import (
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/uptrace/bun"
)

type IdempotencyKey interface {
	data.ICRUDStore[entity.IdempotencyKey]
}

func NewIdempotencyKey(db bun.IDB) IdempotencyKey {
	return data.New[entity.IdempotencyKey](db)
}

func IdemKeyWithKey(key string) data.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("idempotency_key = ?", key)
	}
}

func IdemKeyWithUserID(uid int64) data.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("user_id = ?", uid)
	}
}

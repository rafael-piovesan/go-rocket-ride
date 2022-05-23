package datastore

import (
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/uptrace/bun"
)

type Ride interface {
	data.ICRUDStore[entity.Ride]
}

func NewRide(db bun.IDB) Ride {
	return data.New[entity.Ride](db)
}

func RideWithIdemKeyID(kid int64) data.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("idempotency_key_id = ?", kid)
	}
}

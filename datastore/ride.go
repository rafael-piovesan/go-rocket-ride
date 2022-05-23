package datastore

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/repo"
	"github.com/uptrace/bun"
)

type Ride interface {
	FindAll(context.Context, ...repo.SelectCriteria) ([]entity.Ride, error)
	FindOne(context.Context, ...repo.SelectCriteria) (entity.Ride, error)
	Delete(context.Context, *entity.Ride) error
	Save(context.Context, *entity.Ride) error
	Update(context.Context, *entity.Ride) error
}

type rideStore struct {
	*repo.CRUDRepository[entity.Ride]
}

func NewRide(db bun.IDB) Ride {
	return rideStore{&repo.CRUDRepository[entity.Ride]{DB: db}}
}

func RideWithIdemKeyID(kid int64) repo.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("idempotency_key_id = ?", kid)
	}
}

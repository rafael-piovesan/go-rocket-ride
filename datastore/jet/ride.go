package jet

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rafael-piovesan/go-rocket-ride/entity"

	// dot import so go code would resemble as much as native SQL
	// dot import is not mandatory
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/rafael-piovesan/go-rocket-ride/datastore/jet/rides/public/model"
	. "github.com/rafael-piovesan/go-rocket-ride/datastore/jet/rides/public/table"
)

func (s *sqlStore) CreateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	newRecord := model.Rides{
		ID:               rd.ID,
		CreatedAt:        rd.CreatedAt,
		IdempotencyKeyID: rd.IdempotencyKeyID,
		OriginLat:        rd.OriginLat,
		OriginLon:        rd.OriginLon,
		TargetLat:        rd.TargetLat,
		TargetLon:        rd.TargetLon,
		StripeChargeID:   rd.StripeChargeID,
		UserID:           rd.UserID,
	}

	insertStmt := Rides.INSERT(Rides.AllColumns).
		MODEL(newRecord).
		RETURNING(Rides.AllColumns)

	dest := model.Rides{}

	err := insertStmt.Query(s.db, &dest)

	if err != nil {
		return nil, err
	}

	return toRideEntity(&dest, rd), nil
}

func (s *sqlStore) GetRideByIdempotencyKeyID(ctx context.Context, keyID int64) (*entity.Ride, error) {
	getStmt := Rides.SELECT(Rides.AllColumns).
		FROM(Rides).
		WHERE(Rides.IdempotencyKeyID.EQ(Int(keyID)))

	dest := model.Rides{}

	err := getStmt.Query(s.db, &dest)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return toRideEntity(&dest, nil), nil
}

func (s *sqlStore) UpdateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	updateStmt := Rides.
		UPDATE(Rides.MutableColumns).
		MODEL(Rides).
		WHERE(Rides.ID.EQ(Int(rd.ID))).
		RETURNING(Rides.AllColumns)

	dest := model.Rides{}

	err := updateStmt.Query(s.db, &dest)

	if err != nil {
		return nil, err
	}

	return toRideEntity(&dest, rd), nil
}

func toRideEntity(model *model.Rides, ent *entity.Ride) *entity.Ride {
	if ent == nil {
		ent = &entity.Ride{}
	}

	ent.ID = model.ID
	ent.CreatedAt = model.CreatedAt
	ent.IdempotencyKeyID = model.IdempotencyKeyID
	ent.OriginLat = model.OriginLat
	ent.OriginLon = model.OriginLon
	ent.TargetLat = model.TargetLat
	ent.TargetLon = model.TargetLon
	ent.StripeChargeID = model.StripeChargeID
	ent.UserID = model.UserID

	return ent
}

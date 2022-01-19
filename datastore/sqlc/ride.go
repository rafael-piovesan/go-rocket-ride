package sqlc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

func (s *sqlStore) CreateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	keyID := sql.NullInt64{}
	if rd.IdempotencyKeyID != nil {
		keyID.Int64 = *rd.IdempotencyKeyID
		keyID.Valid = true
	}

	scID := sql.NullString{}
	if rd.StripeChargeID != nil {
		scID.String = *rd.StripeChargeID
		scID.Valid = true
	}

	arg := CreateRideParams{
		CreatedAt:        rd.CreatedAt,
		IdempotencyKeyID: keyID,
		OriginLat:        rd.OriginLat,
		OriginLon:        rd.OriginLon,
		TargetLat:        rd.TargetLat,
		TargetLon:        rd.TargetLon,
		StripeChargeID:   scID,
		UserID:           rd.UserID,
	}
	model, err := s.q.CreateRide(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toRideEntity(&model, rd), nil
}

func (s *sqlStore) GetRideByIdempotencyKeyID(ctx context.Context, keyID int64) (*entity.Ride, error) {
	model, err := s.q.GetRideByIdempotencyKeyID(ctx, sql.NullInt64{Int64: keyID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return toRideEntity(&model, nil), nil
}

func (s *sqlStore) UpdateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	scID := sql.NullString{}
	if rd.StripeChargeID != nil {
		scID.String = *rd.StripeChargeID
		scID.Valid = true
	}

	arg := UpdateRideParams{
		ID:             rd.ID,
		StripeChargeID: scID,
	}

	model, err := s.q.UpdateRide(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toRideEntity(&model, rd), nil
}

func toRideEntity(model *Ride, ent *entity.Ride) *entity.Ride {
	if ent == nil {
		ent = &entity.Ride{}
	}
	var keyID *int64
	if model.IdempotencyKeyID.Valid {
		keyID = &model.IdempotencyKeyID.Int64
	}

	var scID *string
	if model.StripeChargeID.Valid {
		scID = &model.StripeChargeID.String
	}

	ent.ID = model.ID
	ent.CreatedAt = model.CreatedAt
	ent.IdempotencyKeyID = keyID
	ent.OriginLat = model.OriginLat
	ent.OriginLon = model.OriginLon
	ent.TargetLat = model.TargetLat
	ent.TargetLon = model.TargetLon
	ent.StripeChargeID = scID
	ent.UserID = model.UserID

	return ent
}

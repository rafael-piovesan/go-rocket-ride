package bun

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

func (s *sqlStore) CreateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	_, err := s.db.NewInsert().
		Model(rd).
		Returning("*").
		Exec(ctx)

	return rd, err
}

func (s *sqlStore) GetRideByIdempotencyKeyID(ctx context.Context, keyID int64) (*entity.Ride, error) {
	r := &entity.Ride{}
	err := s.db.NewSelect().
		Model(r).
		Where("idempotency_key_id = ?", keyID).
		Limit(1).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return r, entity.ErrNotFound
	}

	return r, err
}

func (s *sqlStore) UpdateRide(ctx context.Context, rd *entity.Ride) (*entity.Ride, error) {
	_, err := s.db.NewUpdate().
		Model(rd).
		WherePK().
		Returning("*").
		Exec(ctx)

	return rd, err
}

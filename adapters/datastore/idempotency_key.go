package datastore

import (
	"context"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

func (s *sqlStore) CreateIdempotencyKey(
	ctx context.Context,
	ik *entity.IdempotencyKey,
) (*entity.IdempotencyKey, error) {
	_, err := s.db.NewInsert().
		Model(ik).
		Returning("*").
		Exec(ctx)

	return ik, err
}

func (s *sqlStore) GetIdempotencyKey(ctx context.Context, key string, userID int64) (*entity.IdempotencyKey, error) {
	ik := entity.IdempotencyKey{}
	err := s.db.NewSelect().
		Model(&ik).
		Where("idempotency_key = ? AND user_id = ?", key, userID).
		Limit(1).
		Scan(ctx)

	return &ik, err
}

func (s *sqlStore) UpdateIdempotencyKey(
	ctx context.Context,
	ik *entity.IdempotencyKey,
) (*entity.IdempotencyKey, error) {
	_, err := s.db.NewUpdate().
		Model(ik).
		WherePK().
		Returning("*").
		Exec(ctx)

	return ik, err
}

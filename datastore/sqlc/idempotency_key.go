package sqlc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
	"github.com/rafael-piovesan/go-rocket-ride/entity/idempotency"
	"github.com/tabbed/pqtype"
)

func (s *sqlStore) CreateIdempotencyKey(
	ctx context.Context,
	ik *entity.IdempotencyKey,
) (*entity.IdempotencyKey, error) {
	lockedAt := sql.NullTime{}
	if ik.LockedAt != nil {
		lockedAt.Time = *ik.LockedAt
		lockedAt.Valid = true
	}

	rCode := sql.NullInt32{}
	if ik.ResponseCode != nil {
		rCode.Int32 = int32(*ik.ResponseCode)
		rCode.Valid = true
	}

	rBody := pqtype.NullRawMessage{}
	if ik.ResponseBody != nil {
		j, err := ik.ResponseBody.Marshal()
		if err != nil {
			return nil, err
		}

		rBody.RawMessage = j
		rBody.Valid = true
	}

	arg := CreateIdempotencyKeyParams{
		CreatedAt:      ik.CreatedAt,
		IdempotencyKey: ik.IdempotencyKey,
		LastRunAt:      ik.LastRunAt,
		LockedAt:       lockedAt,
		RequestMethod:  ik.RequestMethod,
		RequestParams:  ik.RequestParams,
		RequestPath:    ik.RequestPath,
		ResponseCode:   rCode,
		ResponseBody:   rBody,
		RecoveryPoint:  ik.RecoveryPoint.String(),
		UserID:         ik.UserID,
	}

	model, err := s.q.CreateIdempotencyKey(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toIdempotencyKeyEntity(&model, ik)
}

func (s *sqlStore) GetIdempotencyKey(ctx context.Context, key string, userID int64) (*entity.IdempotencyKey, error) {
	arg := GetIdempotencyKeyParams{
		UserID:         userID,
		IdempotencyKey: key,
	}
	model, err := s.q.GetIdempotencyKey(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return toIdempotencyKeyEntity(&model, nil)
}

func (s *sqlStore) UpdateIdempotencyKey(
	ctx context.Context,
	ik *entity.IdempotencyKey,
) (*entity.IdempotencyKey, error) {
	lockedAt := sql.NullTime{}
	if ik.LockedAt != nil {
		lockedAt.Time = *ik.LockedAt
		lockedAt.Valid = true
	}

	rCode := sql.NullInt32{}
	if ik.ResponseCode != nil {
		rCode.Int32 = int32(*ik.ResponseCode)
		rCode.Valid = true
	}

	rBody := pqtype.NullRawMessage{}
	if ik.ResponseBody != nil {
		j, err := ik.ResponseBody.Marshal()
		if err != nil {
			return nil, err
		}

		rBody.RawMessage = j
		rBody.Valid = true
	}

	arg := UpdateIdempotencyKeyParams{
		ID:            ik.ID,
		LastRunAt:     ik.LastRunAt,
		LockedAt:      lockedAt,
		ResponseCode:  rCode,
		ResponseBody:  rBody,
		RecoveryPoint: ik.RecoveryPoint.String(),
	}
	model, err := s.q.UpdateIdempotencyKey(ctx, arg)
	if err != nil {
		return nil, err
	}

	return toIdempotencyKeyEntity(&model, ik)
}

func toIdempotencyKeyEntity(model *IdempotencyKey, ent *entity.IdempotencyKey) (*entity.IdempotencyKey, error) {
	if ent == nil {
		ent = &entity.IdempotencyKey{}
	}
	var lockedAt *time.Time
	if model.LockedAt.Valid {
		lockedAt = &model.LockedAt.Time
	}

	var rCode *int32
	if model.ResponseCode.Valid {
		rCode = &model.ResponseCode.Int32
	}

	var rBody *idempotency.ResponseBody
	if model.ResponseBody.Valid {
		err := json.Unmarshal(model.ResponseBody.RawMessage, &rBody)
		if err != nil {
			return nil, err
		}
	}

	ent.ID = model.ID
	ent.CreatedAt = model.CreatedAt
	ent.IdempotencyKey = model.IdempotencyKey
	ent.LastRunAt = model.LastRunAt
	ent.LockedAt = lockedAt
	ent.RequestMethod = model.RequestMethod
	ent.RequestParams = model.RequestParams
	ent.RequestPath = model.RequestPath
	ent.ResponseCode = (*idempotency.ResponseCode)(rCode)
	ent.ResponseBody = rBody
	ent.RecoveryPoint = idempotency.RecoveryPoint(model.RecoveryPoint)
	ent.UserID = model.UserID

	return ent, nil
}

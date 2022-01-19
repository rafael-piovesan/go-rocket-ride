package sqlc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

func (s *sqlStore) GetUserByEmail(ctx context.Context, e string) (*entity.User, error) {
	model, err := s.q.GetUserByEmail(ctx, e)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return toUserEntity(&model, nil), nil
}

func toUserEntity(model *User, ent *entity.User) *entity.User {
	if ent == nil {
		ent = &entity.User{}
	}

	ent.ID = model.ID
	ent.Email = model.Email
	ent.StripeCustomerID = model.StripeCustomerID

	return ent
}

package datastore

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

func (s *sqlStore) GetUserByEmail(ctx context.Context, e string) (*entity.User, error) {
	user := &entity.User{}
	err := s.db.NewSelect().
		Model(user).
		Where("email = ?", strings.ToLower(e)).
		Limit(1).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrNotFound
	}

	return user, err
}

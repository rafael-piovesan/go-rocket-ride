package jet

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rafael-piovesan/go-rocket-ride/entity"

	. "github.com/go-jet/jet/v2/postgres"
	// dot import so go code would resemble as much as native SQL
	// dot import is not mandatory
	"github.com/rafael-piovesan/go-rocket-ride/datastore/jet/rides/public/model"
	. "github.com/rafael-piovesan/go-rocket-ride/datastore/jet/rides/public/table"
)

func (s *sqlStore) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	getStmt := Users.SELECT(Users.AllColumns).
		FROM(Users).
		WHERE(Users.Email.EQ(String(email)))

	dest := model.Users{}

	err := getStmt.Query(s.db, &dest)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return toUserEntity(&dest, nil), nil
}

func toUserEntity(model *model.Users, ent *entity.User) *entity.User {
	if ent == nil {
		ent = &entity.User{}
	}

	ent.ID = model.ID
	ent.Email = model.Email
	ent.StripeCustomerID = model.StripeCustomerID

	return ent
}

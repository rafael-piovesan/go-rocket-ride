package datastore

import (
	"context"
	"strings"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/repo"
	"github.com/uptrace/bun"
)

type User interface {
	FindAll(context.Context, ...repo.SelectCriteria) ([]entity.User, error)
	FindOne(context.Context, ...repo.SelectCriteria) (entity.User, error)
	Delete(context.Context, *entity.User) error
	Save(context.Context, *entity.User) error
	Update(context.Context, *entity.User) error
}

type userStore struct {
	*repo.CRUDRepository[entity.User]
}

func NewUser(db bun.IDB) User {
	return userStore{&repo.CRUDRepository[entity.User]{DB: db}}
}

func UserWithEmail(e string) repo.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("email = ?", strings.ToLower(e))
	}
}

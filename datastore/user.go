package datastore

import (
	"strings"

	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
	"github.com/rafael-piovesan/go-rocket-ride/v2/pkg/data"
	"github.com/uptrace/bun"
)

type User interface {
	data.ICRUDStore[entity.User]
}

func NewUser(db bun.IDB) User {
	return data.New[entity.User](db)
}

func UserWithEmail(e string) data.SelectCriteria {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("email = ?", strings.ToLower(e))
	}
}

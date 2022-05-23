package context

import (
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

func AddUser(c echo.Context, u entity.User) {
	c.Set("my-user", u)
}

func GetUser(c echo.Context) (ur entity.User, ok bool) {
	ur, ok = c.Get("my-user").(entity.User)
	return
}

func AddIdemKey(c echo.Context, ik entity.IdempotencyKey) {
	c.Set("my-key", ik)
}

func GetIdemKey(c echo.Context) (ik entity.IdempotencyKey, ok bool) {
	ik, ok = c.Get("my-key").(entity.IdempotencyKey)
	return
}

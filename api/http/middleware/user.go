package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/v2/datastore"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

type userRequest struct {
	UserKey string `header:"authorization" validate:"required"`
}

func NewUser(store datastore.User) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			binder := &echo.DefaultBinder{}
			validate := validator.New()
			ur := userRequest{}

			if err := binder.BindHeaders(c, &ur); err != nil {
				return err
			}

			if err := validate.Struct(&ur); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			// This is obviously something you shouldn't do in a real application, but for
			// now we're just going to trust that the user is whoever they said they were
			// from an email in the `Authorization` header.
			user, err := store.FindOne(c.Request().Context(), datastore.UserWithEmail(ur.UserKey))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, entity.ErrPermissionDenied.Error())
			}

			c.Set(string(entity.UserCtxKey), user)

			return next(c)
		}
	}
}

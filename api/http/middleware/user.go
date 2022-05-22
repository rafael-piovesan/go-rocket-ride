package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	rocketride "github.com/rafael-piovesan/go-rocket-ride/v2"
	"github.com/rafael-piovesan/go-rocket-ride/v2/entity"
)

type userRequest struct {
	UserKey string `header:"authorization" validate:"required"`
}

type UserMiddleware struct {
	binder   *echo.DefaultBinder
	validate *validator.Validate
	store    rocketride.Datastore
}

func NewUserMiddleware(ds rocketride.Datastore) *UserMiddleware {
	return &UserMiddleware{
		binder:   &echo.DefaultBinder{},
		validate: validator.New(),
		store:    ds,
	}
}

func (u *UserMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ur := userRequest{}
		if err := u.binder.BindHeaders(c, &ur); err != nil {
			return err
		}

		if err := u.validate.Struct(&ur); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// This is obviously something you shouldn't do in a real application, but for
		// now we're just going to trust that the user is whoever they said they were
		// from an email in the `Authorization` header.
		user, err := u.store.GetUserByEmail(c.Request().Context(), ur.UserKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, entity.ErrPermissionDenied.Error())
		}

		c.Set(string(entity.UserCtxKey), user)

		return next(c)
	}
}

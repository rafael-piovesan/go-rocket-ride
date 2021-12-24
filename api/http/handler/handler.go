package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rafael-piovesan/go-rocket-ride/entity"
)

type Handler struct {
	binder   *echo.DefaultBinder
	validate *validator.Validate
}

func New() *Handler {
	return &Handler{
		binder:   &echo.DefaultBinder{},
		validate: validator.New(),
	}
}

// BindAndValidate binds and validates the data struct pointed to by 'i'.
// It expects a stuct pointer as parameter, since it needs to populate its
// fields, otherwise it'll panic.
func (h *Handler) BindAndValidate(c echo.Context, i interface{}) error {
	if i == nil {
		return nil
	}

	if err := h.binder.BindHeaders(c, i); err != nil {
		return err
	}

	if err := h.binder.BindBody(c, i); err != nil {
		return err
	}

	if err := h.validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (h *Handler) GetUserID(c echo.Context) (uid int64, err error) {
	uid, ok := c.Get("user-id").(int64)
	if !ok {
		err = echo.NewHTTPError(http.StatusUnauthorized, entity.ErrPermissionDenied.Error())
	}
	return
}
